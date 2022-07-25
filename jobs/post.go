package jobs

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/go-playground/validator.v9"
)

// swagger:route POST /jobs Jobs post-job
//
// Add Job
//
// You can add a Job through this endpoint by filling in the details of the job to be added and the organization which is posting the job.
//
// responses:
//  201: Jobs
//  400: Error

func PostJob(w http.ResponseWriter, r *http.Request) {

	var add structures.Jobs

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	id := uuid.New()
	add.ID = id.String()
	add.Org_id = dataToCompare["org_id"]
	add.Cat_ID = dataToCompare["cat_id"]
	add.Designation = dataToCompare["designation"]
	add.Description = dataToCompare["description"]
	add.Location = dataToCompare["location"]
	add.Salary = dataToCompare["salary"]

	// input validation
	validate := validator.New()
	err := validate.Struct(add)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	// validating json schema
	schemaLoader := gojsonschema.NewReferenceLoader("file:///home/sarmad/Go_Practice/PostJson/structures/JobSchema.json")
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

	result := db.Conn.Create(&add)
	if result.Error != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not enter record!!")
		return
	} else {
		w.WriteHeader(201)
		fmt.Fprintf(w, "Job Created!!")
	}
}
