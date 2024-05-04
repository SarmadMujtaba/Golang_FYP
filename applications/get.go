package applications

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

type Resp struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Org_id      string `json:"org_id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Designation string `json:"designation"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Salary      string `json:"salary"`
	Status      string `json:"status"`
	// CreatedAt   time.Time `json:"CreatedAt"`
}

func GetApplications(w http.ResponseWriter, r *http.Request) {
	var app structures.Applications
	// var apps []structures.Applications
	// var jobs []structures.Jobs
	var resp []Resp
	isEmpty := true

	app.U_ID = strings.ReplaceAll(r.URL.Query().Get("user_id"), `"`, "")
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

		result := db.Conn.
			Table("jobs").
			Joins("JOIN organizations ON jobs.org_id = organizations.org_id").
			Joins("JOIN applications ON jobs.id = applications.job_id").
			Where("applications.u_id = ?", app.U_ID).
			Select("jobs.id, jobs.org_id, jobs.designation, jobs.description, jobs.location, jobs.salary, organizations.name AS name, applications.status").
			Scan(&resp)

		// result := db.Conn.Where("U_ID = ?", app.U_ID).Find(&app)
		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return")
			return
		}
		// userID := app.U_ID
		// db.Conn.Joins("JOIN jobs ON applications.job_id = jobs.id").Select("applications.u_id, applications.status, jobs.*").Where("applications.u_id = ?", userID).Find(&apps).Scan(&jobs)

		// var response []Response

		// for i, job := range jobs {
		// 	response = append(response, Response{
		// 		U_ID:   apps[i].U_ID,
		// 		Status: apps[i].Status,
		// 		Job:    job,
		// 	})
		// }
		json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		isEmpty = false
		return
	}

	app.Job_ID = strings.ReplaceAll(r.URL.Query().Get("job_id"), `"`, "")
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
