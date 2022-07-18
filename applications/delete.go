package applications

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

func DeleteApplications(w http.ResponseWriter, r *http.Request) {
	var app structures.Applications
	var apps []structures.Applications
	wrongInput := true

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

		db.Conn.Where("U_ID = ?", app.U_ID).Find(&apps)
		for _, row := range apps {
			if row.U_ID == app.U_ID {
				db.Conn.Where("U_ID = ?", app.U_ID).Delete(&apps)
				w.WriteHeader(200)
				fmt.Fprintf(w, "Application Deleted!!")
				wrongInput = false
				return
			}
		}
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

		db.Conn.Where("Job_ID = ?", app.Job_ID).Find(&apps)
		for _, row := range apps {
			if row.Job_ID == app.Job_ID {
				db.Conn.Where("Job_ID = ?", app.Job_ID).Delete(&apps)
				w.WriteHeader(200)
				fmt.Fprintf(w, "Application Deleted!!")
				wrongInput = false
				return
			}
		}
	}

	if wrongInput {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not delete application!!")
		return
	}

}
