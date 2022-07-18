package applications

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

func GetApplications(w http.ResponseWriter, r *http.Request) {
	var app structures.Applications
	var apps []structures.Applications
	isEmpty := true

	app.U_ID = r.URL.Query().Get("user_id")
	if len(app.U_ID) > 0 {
		// populating add for validation
		app.Job_ID = app.U_ID
		validate := validator.New()
		err := validate.Struct(app)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		result := db.Conn.Where("U_ID = ?", app.U_ID).Find(&apps)
		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return")
			return
		}

		if len(apps) == 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return")
			return
		}

		json.Marshal(apps)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(apps)
		isEmpty = false
		return
	}

	app.Job_ID = r.URL.Query().Get("job_id")
	if len(app.Job_ID) > 0 {
		// populating add for validation
		app.U_ID = app.Job_ID
		validate := validator.New()
		err := validate.Struct(app)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		result := db.Conn.Where("Job_ID = ?", app.Job_ID).Find(&apps)
		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return")
			return
		}

		if len(apps) == 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return")
			return
		}

		json.Marshal(apps)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(apps)
		isEmpty = false
		return
	}

	if isEmpty {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Missing or wrong Parameters!!")
	}

}
