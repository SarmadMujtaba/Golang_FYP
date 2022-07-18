package jobs

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

func AddSkill(w http.ResponseWriter, r *http.Request) {
	var skill structures.RequiredSkills
	var jobs []structures.Jobs

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	skill.ID = dataToCompare["job_id"]
	skill.Skill = dataToCompare["skill"]

	validate := validator.New()
	err := validate.Struct(skill)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect Input")
		return
	}

	db.Conn.Where("ID = ?", skill.ID).Find(&jobs)
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
