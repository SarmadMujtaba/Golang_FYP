package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

type Users struct {
	ID    string `json:"id" validate:"uuid"`
	Name  string `json:"name" validate:"alpha"`
	Email string `json:"email" validate:"email"`
	Pass  string `json:"pass" validate:"alphanum"`
}

type Organizations struct {
	U_ID    string `json:"user_id" validate:"uuid"`
	Org_ID  string `json:"id" validate:"uuid"`
	Name    string `json:"name" validate:"alpha"`
	About   string `json:"about"`
	Website string `json:"website" validate:"url"`
}

var conn *sql.DB

func main() {
	const (
		host     = "localhost"
		port     = 3306
		user     = "root"
		password = "DummySQL786"
		dbname   = "db"
	)
	connString := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", user, password, host, port, dbname)

	var err error
	conn, err = sql.Open("mysql", connString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connection Established...")
	Handler()
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []Users
	var user Users
	var add Users

	add.ID = r.URL.Query().Get("id")
	if len(add.ID) > 0 {
		// populating add for validation
		add.Name = "dummy"
		add.Pass = "dummy"
		add.Email = "dummy@gmail.com"
		validate := validator.New()
		err := validate.Struct(add)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		err2 := conn.QueryRow("SELECT * FROM user_data where id = ?", add.ID).Scan(&user.ID, &user.Email, &user.Name, &user.Pass)
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
		w.WriteHeader(404)
		fmt.Fprintf(w, "Failed to Fetch!!")
		return
	}
	for getAll.Next() {
		getAll.Scan(&user.ID, &user.Email, &user.Name, &user.Pass)
		users = append(users, user)
	}
	if len(users) == 0 {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Nothing to return!!")
		return
	}

	json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	return
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	var add Users

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	id := uuid.New()
	add.ID = id.String()
	add.Email = dataToCompare["email"]
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
	conn.QueryRow("Select name from user_data where email = ?", add.Email).Scan(&duplicate)
	if duplicate != "" {
		w.WriteHeader(409)
		fmt.Fprintf(w, "Email ID already exist!!")
		return
	}

	_, err = conn.Query("insert into user_data values(?, ?, ?, ?)", add.ID, add.Email, add.Name, add.Pass)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(201)
	fmt.Fprintf(w, "User inserted!!")
	return
}

func deleteUsers(w http.ResponseWriter, r *http.Request) {
	var check Users
	check.ID = r.URL.Query().Get("id")
	if len(check.ID) > 0 {
		// populating add for validation
		check.Name = "dummy"
		check.Email = "dummy@gmail.com"
		check.Pass = "dummy"
		validate := validator.New()
		err := validate.Struct(check)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}
	}

	var find string
	conn.QueryRow("Select name from user_data where id = ?", check.ID).Scan(&find)
	if find == "" {
		w.WriteHeader(404)
		fmt.Fprintf(w, "User does not exist!!")
		return
	}
	var org string
	conn.QueryRow("select org_id from membership where id = ?", check.ID).Scan(&org)
	conn.Query("DELETE from membership where id = ?", check.ID)
	conn.Query("DELETE from organizations where org_id = ?", org)
	conn.Query("DELETE from user_data where id = ?", check.ID)
	w.WriteHeader(200)
	fmt.Fprintf(w, "Record deleted successfully!!")
	return
}

func postOrganizations(w http.ResponseWriter, r *http.Request) {
	var add Organizations

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	id := uuid.New()
	add.Org_ID = id.String()
	add.U_ID = dataToCompare["user_id"]
	add.Name = dataToCompare["name"]
	add.About = dataToCompare["about"]
	add.Website = dataToCompare["website"]

	// input validation
	validate := validator.New()
	err := validate.Struct(add)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	var duplicate string
	conn.QueryRow("Select name from organizations where name = ?", add.Name).Scan(&duplicate)
	if duplicate != "" {
		w.WriteHeader(409)
		fmt.Fprintf(w, "Organization already exist!!")
		return
	}

	_, er := conn.Query("insert into organizations values(?, ?, ?, ?)", add.Org_ID, add.Name, add.About, add.Website)
	if er != nil {
		fmt.Fprintln(w, "Could not enter record!!")
		return
	}

	id = uuid.New()
	_, err = conn.Query("insert into membership values(?, ?, ?)", id.String(), add.U_ID, add.Org_ID)
	if err != nil {
		w.WriteHeader(400)
		conn.Query("delete from organizations where org_id = ?", add.Org_ID)
		fmt.Fprintln(w, "Could not enter record!!")
		return
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, "Organization Created!!")
	return
}

func getOrganizations(w http.ResponseWriter, r *http.Request) {
	var organizations []Organizations
	var organization Organizations
	var add Organizations

	add.Org_ID = r.URL.Query().Get("id")
	if len(add.Org_ID) > 0 {
		// populating add for validation
		add.U_ID = "73c6ba9b-9325-4c68-bacb-52b6ce04e919"
		add.Name = "dummy"
		add.About = "dummy"
		add.Website = "https://pkg.go.dev/github.com"
		validate := validator.New()
		err := validate.Struct(add)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		err2 := conn.QueryRow("SELECT org_id, name, about, website FROM organizations where org_id = ?", add.Org_ID).Scan(&organization.Org_ID, &organization.Name, &organization.About, &organization.Website)
		err3 := conn.QueryRow("select id from membership where org_id = ?", add.Org_ID).Scan(&organization.U_ID)
		if err2 != nil || err3 != nil {
			w.WriteHeader(404)
			fmt.Fprintf(w, "Organization Not Found!!")
			return
		}
		json.Marshal(organization)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(organization)
		return
	}

	getAll, err := conn.Query("select * from organizations")
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Record not found!!")
		return
	}
	for getAll.Next() {
		getAll.Scan(&organization.Org_ID, &organization.Name, &organization.About, &organization.Website)
		conn.QueryRow("select id from membership where org_id = ?", organization.Org_ID).Scan(&organization.U_ID)
		organizations = append(organizations, organization)
	}

	if len(organizations) == 0 {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Nothing to return!!")
		return
	}

	json.Marshal(organizations)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(organizations)
	return
}

func deleteOrganizations(w http.ResponseWriter, r *http.Request) {
	var check Organizations
	check.Org_ID = r.URL.Query().Get("id")
	if len(check.Org_ID) > 0 {
		// populating add for validation
		check.U_ID = "73c6ba9b-9325-4c68-bacb-52b6ce04e919"
		check.Name = "dummy"
		check.About = "dummy"
		check.Website = "https://pkg.go.dev/github.com"
		validate := validator.New()
		err := validate.Struct(check)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}
	}

	var find string
	conn.QueryRow("Select name from organizations where org_id = ?", check.Org_ID).Scan(&find)
	if find == "" {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Organization does not exist!!")
		return
	}
	conn.Query("DELETE from membership where org_id = ?", check.Org_ID)
	conn.Query("DELETE from organizations where org_id = ?", check.Org_ID)
	w.WriteHeader(200)
	fmt.Fprintf(w, "Record deleted successfully!!")
	return
}

func Handler() {
	route := mux.NewRouter()
	route.HandleFunc("/users", getUsers).Methods(http.MethodGet)
	route.HandleFunc("/users", postUsers).Methods(http.MethodPost)
	route.HandleFunc("/users", deleteUsers).Methods(http.MethodDelete)
	route.HandleFunc("/organizations", getOrganizations).Methods(http.MethodGet)
	route.HandleFunc("/organizations", postOrganizations).Methods(http.MethodPost)
	route.HandleFunc("/organizations", deleteOrganizations).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":5020", route))
}
