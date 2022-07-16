package jobs

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

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
