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
)

type Users struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Pass string `json:"pass"`
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

	insert, err := conn.Query("select * from user_data")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to Fetch!!")
	}
	for insert.Next() {
		insert.Scan(&user.ID, &user.Name, &user.Pass)
		users = append(users, user)
	}
	var count int = 0
	id := r.URL.Query().Get("id")
	if len(id) > 0 {
		for _, v := range users {
			if v.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(users[count])
				return
			}
			count++
		}
	}
	json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	return
}

func postData(w http.ResponseWriter, r *http.Request) {
	var users []Users
	var user Users
	var add Users

	get, err := conn.Query("select * from user_data")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to Insert!!")
	}
	for get.Next() {
		get.Scan(&user.ID, &user.Name, &user.Pass)
		users = append(users, user)
	}
	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)
	for _, v := range users {
		if v.ID == dataToCompare["id"] {
			w.WriteHeader(409)
			fmt.Fprintf(w, "ID already in use!!")
			return
		}
	}
	add.ID = dataToCompare["id"]
	add.Name = dataToCompare["name"]
	add.Pass = dataToCompare["pass"]
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
