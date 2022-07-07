package main

import (
	"PostJson/db"
	"PostJson/organizations"
	"PostJson/structures"
	"PostJson/users"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db.Connection()
	db.Conn.AutoMigrate(&structures.Users{}, &structures.Organizations{}, &structures.Memberships{})
	Handler()
}

func Handler() {
	route := mux.NewRouter()
	route.HandleFunc("/users", users.GetUsers).Methods(http.MethodGet)
	route.HandleFunc("/users", users.PostUsers).Methods(http.MethodPost)
	route.HandleFunc("/users", users.DeleteUsers).Methods(http.MethodDelete)
	route.HandleFunc("/organizations", organizations.GetOrganizations).Methods(http.MethodGet)
	route.HandleFunc("/organizations", organizations.PostOrganizations).Methods(http.MethodPost)
	route.HandleFunc("/organizations", organizations.DeleteOrganizations).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":5020", route))
}
