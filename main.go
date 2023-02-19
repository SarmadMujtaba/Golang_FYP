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
	"PostJson/authentication"
	"PostJson/db"
	"PostJson/invites"
	"PostJson/jobs"
	"PostJson/login"
	"PostJson/members"
	"PostJson/organizations"
	userprofile "PostJson/userProfile"
	"PostJson/users"
	"log"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Loading Environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
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

	route.HandleFunc("/login", login.Login).Methods(http.MethodPost)
	route.HandleFunc("/users", authentication.IsAuthorized(users.GetUsers)).Methods(http.MethodGet)
	route.HandleFunc("/signup", authentication.VerifyEmail(login.Signup)).Methods(http.MethodPost)
	route.HandleFunc("/users", authentication.IsAuthorized(users.DeleteUsers)).Methods(http.MethodDelete)
	route.HandleFunc("/organizations", authentication.IsAuthorized(organizations.GetOrganizations)).Methods(http.MethodGet)
	route.HandleFunc("/organizations", authentication.IsAuthorized(organizations.PostOrganizations)).Methods(http.MethodPost)
	route.HandleFunc("/organizations", authentication.IsAuthorized(organizations.DeleteOrganizations)).Methods(http.MethodDelete)
	route.HandleFunc("/invite", authentication.IsAuthorized(invites.GetInvites)).Methods(http.MethodGet)
	route.HandleFunc("/invite", authentication.IsAuthorized(invites.PostInvite)).Methods(http.MethodPost)
	route.HandleFunc("/members", authentication.IsAuthorized(members.GetMembers)).Methods(http.MethodGet)
	route.HandleFunc("/members", authentication.IsAuthorized(members.PostMembers)).Methods(http.MethodPost)
	route.HandleFunc("/jobs", authentication.IsAuthorized(jobs.GetJobs)).Methods(http.MethodGet)
	route.HandleFunc("/jobs/designations", authentication.IsAuthorized(jobs.GetDesignations)).Methods(http.MethodGet)
	route.HandleFunc("/jobs", authentication.IsAuthorized(jobs.PostJob)).Methods(http.MethodPost)
	route.HandleFunc("/jobs", authentication.IsAuthorized(jobs.DeleteJob)).Methods(http.MethodDelete)
	route.HandleFunc("/jobs/skills", authentication.IsAuthorized(jobs.AddSkill)).Methods(http.MethodPost)
	route.HandleFunc("/category", authentication.IsAuthorized(jobs.GetCategory)).Methods(http.MethodGet)
	route.HandleFunc("/profile", authentication.IsAuthorized(userprofile.GetProfile)).Methods(http.MethodGet)
	route.HandleFunc("/profile", authentication.IsAuthorized(userprofile.Profile)).Methods(http.MethodPut)
	route.HandleFunc("/profile/skills", authentication.IsAuthorized(userprofile.AddSkill)).Methods(http.MethodPost)
	route.HandleFunc("/profile/experience", authentication.IsAuthorized(userprofile.AddExperience)).Methods(http.MethodPost)
	route.HandleFunc("/application", authentication.IsAuthorized(applications.PostApplication)).Methods(http.MethodPost)
	route.HandleFunc("/application", authentication.IsAuthorized(applications.GetApplications)).Methods(http.MethodGet)
	route.HandleFunc("/application", authentication.IsAuthorized(applications.DeleteApplications)).Methods(http.MethodDelete)
	route.HandleFunc("/verify", authentication.Verify).Methods(http.MethodGet)

	route.HandleFunc("/upload", applications.FileUpload).Methods(http.MethodPost)

	route.HandleFunc("/application/shortlist", applications.Shortlist).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":5020", route))
}
