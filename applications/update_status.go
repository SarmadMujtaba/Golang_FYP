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

func UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var app structures.Applications
	var apps []structures.Applications

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
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	// Update the status
	result := db.Conn.Model(&apps).Where("U_ID = ? AND Job_ID = ?", app.U_ID, app.Job_ID).Update("Status", app.Status)
	if result.Error != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Failed to update status")
		return
	}
	w.WriteHeader(200)
	fmt.Fprintln(w, "Status Updated!!")
}
