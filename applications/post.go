package applications

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// swagger:route POST /application Application post-application
//
// Add Application
//
// You can add a user's application to a job through this endpoint by filling in the details of the user and the job.
//
// responses:
//  200: Applications
//  400: Error
//  409: Error

func PostApplication(w http.ResponseWriter, r *http.Request) {
	// var user structures.Users
	// var org structures.Organizations
	var app structures.Applications
	var repeatCheck []structures.Applications

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	app.U_ID = strings.ReplaceAll(dataToCompare["user_id"], `"`, "")
	app.Job_ID = strings.ReplaceAll(dataToCompare["job_id"], `"`, "")
	app.Status = dataToCompare["status"]

	fmt.Println(app.U_ID)
	fmt.Println(app.Job_ID)
	fmt.Println(app.Status)

	// validate := validator.New()
	// err := validate.Struct(app)
	// if err != nil {
	// 	w.WriteHeader(400)
	// 	fmt.Fprintf(w, "Incorrect Input")
	// 	return
	// }

	// validating json schema
	schemaLoader := gojsonschema.NewReferenceLoader("file:///app/schemas/ApplicationSchema.json")
	documentLoader := gojsonschema.NewGoLoader(dataToCompare)

	res, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}
	if !res.Valid() {
		w.WriteHeader(400)
		for _, desc := range res.Errors() {
			fmt.Fprintln(w, desc.Description())
		}
		return
	}

	db.Conn.Where("U_ID = ?", app.U_ID).Find(&repeatCheck)
	for _, row := range repeatCheck {
		if row.Job_ID == app.Job_ID {
			w.WriteHeader(409)
			fmt.Fprintf(w, "This user already applied for this job!!")
			return
		}
	}

	result := db.Conn.Create(&app)
	if result.Error != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not enter record!!")
		return
	}
	w.WriteHeader(200)
	fmt.Fprintln(w, "Application Submitted!!")
}
