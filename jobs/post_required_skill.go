package jobs

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	var jobs []structures.Jobs

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]interface{}
	json.Unmarshal(dataFromWeb, &dataToCompare)

	jobID := dataToCompare["job_id"].(string)

	var skills []string
	skillsRaw := dataToCompare["skills"].([]interface{})
	for _, s := range skillsRaw {
		skills = append(skills, s.(string))
	}

	fmt.Println(jobID, skills)

	for _, v := range skills {
		skill := structures.RequiredSkills{
			Job_ID: strings.ReplaceAll(jobID, `"`, ""),
			Skill:  v,
		}
		fmt.Println(skill.Job_ID)
		fmt.Println(skill.Skill)

		// validate := validator.New()
		// err := validate.Struct(skill)
		// if err != nil {
		//  w.WriteHeader(400)
		//  fmt.Fprintf(w, "Incorrect Input")
		//  return
		// }

		// validating json schema
		// schemaLoader := gojsonschema.NewReferenceLoader("file:///app/schemas/ReqSkillSchema.json")
		// documentLoader := gojsonschema.NewGoLoader(dataToCompare)

		// res, err := gojsonschema.Validate(schemaLoader, documentLoader)
		// if err != nil {
		// 	panic(err.Error())
		// }
		// if !res.Valid() {
		// 	w.WriteHeader(400)
		// 	for _, desc := range res.Errors() {
		// 		fmt.Fprintln(w, desc.Description())
		// 	}
		// 	return
		// }

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
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, "Required Skills added!!")
}
