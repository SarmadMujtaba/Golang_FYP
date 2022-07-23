//  Package classification Final Year Project - APIs.
//   version: 0.0.1
//   title: Final Year Project
//  Schemes: http
//  Host: localhost:5020
//  BasePath: /
//  Consumes:
//    - application/json
//  Produces:
//    - application/json
//  Contact: Sarmad Mujtaba <srmdmjtba@gmail.com>
//
// swagger:meta
package main

import (
	"PostJson/applications"
	"PostJson/db"
	"PostJson/jobs"
	"PostJson/members"
	"PostJson/organizations"
	userprofile "PostJson/userProfile"
	"PostJson/users"
	"log"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db.Connection()
	Handler()
}

func Handler() {
	route := mux.NewRouter()

	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	route.Handle("/docs", sh)

	route.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	route.HandleFunc("/users", users.GetUsers).Methods(http.MethodGet)
	route.HandleFunc("/users", users.PostUsers).Methods(http.MethodPost)
	route.HandleFunc("/users", users.DeleteUsers).Methods(http.MethodDelete)
	route.HandleFunc("/organizations", organizations.GetOrganizations).Methods(http.MethodGet)
	route.HandleFunc("/organizations", organizations.PostOrganizations).Methods(http.MethodPost)
	route.HandleFunc("/organizations", organizations.DeleteOrganizations).Methods(http.MethodDelete)
	route.HandleFunc("/members", members.GetMembers).Methods(http.MethodGet)
	route.HandleFunc("/members", members.PostMembers).Methods(http.MethodPost)
	route.HandleFunc("/jobs", jobs.GetJobs).Methods(http.MethodGet)
	route.HandleFunc("/jobs/designations", jobs.GetDesignations).Methods(http.MethodGet)
	route.HandleFunc("/jobs", jobs.PostJob).Methods(http.MethodPost)
	route.HandleFunc("/jobs", jobs.DeleteJob).Methods(http.MethodDelete)
	route.HandleFunc("/jobs/skills", jobs.AddSkill).Methods(http.MethodPost)
	route.HandleFunc("/category", jobs.GetCategory).Methods(http.MethodGet)
	route.HandleFunc("/profile", userprofile.GetProfile).Methods(http.MethodGet)
	route.HandleFunc("/profile", userprofile.Profile).Methods(http.MethodPut)
	route.HandleFunc("/profile/skills", userprofile.AddSkill).Methods(http.MethodPost)
	route.HandleFunc("/profile/experience", userprofile.AddExperience).Methods(http.MethodPost)
	route.HandleFunc("/application", applications.PostApplication).Methods(http.MethodPost)
	route.HandleFunc("/application", applications.GetApplications).Methods(http.MethodGet)
	route.HandleFunc("/application", applications.DeleteApplications).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":5020", route))
}
