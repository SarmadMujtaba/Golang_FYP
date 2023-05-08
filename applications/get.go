package applications

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

type Response struct {
	U_ID   string
	Status string
	Job    structures.Jobs
}

type Response2 struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	Education   string `json:"education"`
	Phone       string `json:"phone"`
	Skills      string `json:"skills"`
	Experiences string `json:"experiences"`
}

func GetApplications(w http.ResponseWriter, r *http.Request) {
	var app structures.Applications
	var apps []structures.Applications
	var jobs []structures.Jobs
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

		result := db.Conn.Where("U_ID = ?", app.U_ID).Find(&app)
		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return")
			return
		}
		userID := app.U_ID
		db.Conn.Joins("JOIN jobs ON applications.job_id = jobs.id").Select("applications.u_id, applications.status, jobs.*").Where("applications.u_id = ?", userID).Find(&apps).Scan(&jobs)

		var response []Response

		for i, job := range jobs {
			response = append(response, Response{
				U_ID:   apps[i].U_ID,
				Status: apps[i].Status,
				Job:    job,
			})
		}
		json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
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

		result := db.Conn.Where("Job_ID = ?", app.Job_ID).Find(&app)
		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return")
			return
		}

		var res []Response2

		db.Conn.Table("applications").
			Select("applications.u_id as user_id, users.name, users.email, applications.status, profiles.education, profiles.phone, "+
				"(SELECT GROUP_CONCAT(skills.skill SEPARATOR ', ') FROM skills WHERE skills.u_id = applications.u_id) AS skills, "+
				"(SELECT GROUP_CONCAT(experiences.experience SEPARATOR ', ') FROM experiences WHERE experiences.u_id = applications.u_id) as experiences").
			Joins("LEFT JOIN profiles ON profiles.u_id = applications.u_id").
			Joins("LEFT JOIN users ON users.id = applications.u_id").
			Where("applications.job_id = ?", r.URL.Query().Get("job_id")).
			Scan(&res)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		isEmpty = false
		return
	}

	if isEmpty {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Missing or wrong Parameters!!")
	}

}
