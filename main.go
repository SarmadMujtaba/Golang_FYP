package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

type Users struct {
	Name string `json:"name" validate:"alpha"`
	ID   string `json:"id" validate:"numeric"`
	Pass string `json:"pass" validate:"alphanum"`
}

var conn *sql.DB

func main() {
	const (
		host     = "localhost"
		port     = 3306
		user     = "root"
		password = "DummySQL786"
		dbname   = "users"
	)
	connString := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", user, password, host, port, dbname)

	var err error
	conn, err = sql.Open("mysql", connString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connection Established...")
	API()
}

func getData(w http.ResponseWriter, r *http.Request) {
	var users []Users
	var user Users
	var add Users

	add.ID = r.URL.Query().Get("id")
	if len(add.ID) > 0 {
		add.Name = "test"
		add.Pass = "test"
		validate := validator.New()
		err := validate.Struct(add)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		err2 := conn.QueryRow("SELECT * FROM user_data where id = ?", add.ID).Scan(&user.ID, &user.Name, &user.Pass)
		if err2 != nil {
			w.WriteHeader(404)
			fmt.Fprintf(w, "User Not Found!!")
			return
		}
		json.Marshal(user)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return
	}

	getAll, err := conn.Query("select * from user_data")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to Fetch!!")
	}
	for getAll.Next() {
		getAll.Scan(&user.ID, &user.Name, &user.Pass)
		users = append(users, user)
	}

	json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	return
}

func postData(w http.ResponseWriter, r *http.Request) {
	var add Users

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	ok := json.Valid(dataFromWeb)
	if !ok {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect Syntax!!")
		return
	}
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	add.ID = dataToCompare["id"]
	add.Name = dataToCompare["name"]
	add.Pass = dataToCompare["pass"]

	// input validation
	validate := validator.New()
	err := validate.Struct(add)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	var duplicate string
	conn.QueryRow("Select name from user_data where id = ?", dataToCompare["id"]).Scan(&duplicate)
	if duplicate != "" {
		w.WriteHeader(409)
		fmt.Fprintf(w, "ID already in use!!")
		return
	}

	conn.Query("insert into user_data values(?, ?, ?)", add.ID, add.Name, add.Pass)
	w.WriteHeader(201)
	fmt.Fprintf(w, "User inserted!!")
	return

}

func API() {
	r := mux.NewRouter()
	r.HandleFunc("/users", getData).Methods("GET")
	r.HandleFunc("/users", postData).Methods("POST")
	log.Fatal(http.ListenAndServe(":5010", r))
}
