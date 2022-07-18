package applications

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

func PostApplication(w http.ResponseWriter, r *http.Request) {
	// var user structures.Users
	// var org structures.Organizations
	var app structures.Applications
	var repeatCheck []structures.Applications

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	app.U_ID = dataToCompare["user_id"]
	app.Job_ID = dataToCompare["job_id"]
	app.Status = dataToCompare["status"]

	validate := validator.New()
	err := validate.Struct(app)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect Input")
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
	fmt.Fprintln(w, "Application Submitted!!")
}
