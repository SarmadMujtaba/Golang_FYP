package jobs

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/go-playground/validator.v9"
)

// swagger:route POST /jobs/skills Jobs post-RequiredSkill
//
// Add Required Skill
//
// You can add multiple required skills for a job through this endpoint by filling in the details.
//
// responses:
//  201: RequiredSkills
//  400: Error

func AddSkill(w http.ResponseWriter, r *http.Request) {
	var skill structures.RequiredSkills
	var jobs []structures.Jobs

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	skill.Job_ID = dataToCompare["job_id"]
	skill.Skill = dataToCompare["skill"]

	validate := validator.New()
	err := validate.Struct(skill)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect Input")
		return
	}

	// validating json schema
	schemaLoader := gojsonschema.NewReferenceLoader("file:///home/sarmad/Go_Practice/PostJson/schemas/ReqSkillSchema.json")
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

	db.Conn.Where("ID = ?", skill.Job_ID).Find(&jobs)
	if len(jobs) == 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "job does not exist!!")
		return
	}

	result := db.Conn.Create(&skill)
	if result.Error != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Could not add Required skill!!")
		return
	}
	w.WriteHeader(201)
	fmt.Fprintf(w, "Required Skill added!!")
}
