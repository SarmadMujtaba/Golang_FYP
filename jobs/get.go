package jobs

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Resp struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Org_id      string    `json:"org_id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Designation string    `json:"designation"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Salary      string    `json:"salary"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

func GetJobs(w http.ResponseWriter, r *http.Request) {

	var jobs []structures.Jobs
	var job structures.Jobs
	var resp []Resp

	job.Org_id = strings.ReplaceAll(r.URL.Query().Get("org_id"), `"`, "")

	if len(job.Org_id) > 0 {
		// populating add for validation
		job.ID = job.Org_id
		// validate := validator.New()
		// err := validate.Struct(job)
		// if err != nil {
		// 	w.WriteHeader(400)
		// 	fmt.Fprintf(w, "Incorrect input!!")
		// 	return
		// }
		// Returns all jobs of an organization
		db.Conn.Model(&jobs).
			Joins("JOIN organizations ON jobs.org_id = organizations.org_id").
			Select("jobs.*, organizations.name").
			Where("jobs.org_id = ?", job.Org_id).
			Scan(&resp)

		// db.Conn.Find(&jobs, "Org_id = ?", job.Org_id)

		if len(resp) == 0 {
			w.WriteHeader(200)
			// fmt.Fprintf(w, "Nothing to return!!")
			return
		}

		json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
	// Returns all jobs against a designation name
	job.Designation = strings.ReplaceAll(r.URL.Query().Get("job_name"), `"`, "")
	if len(job.Designation) > 0 {

		db.Conn.Model(&jobs).
			Joins("JOIN organizations ON jobs.org_id = organizations.org_id").
			Select("jobs.*, organizations.name").
			Where("jobs.designation LIKE ?", "%"+job.Designation+"%").
			Scan(&resp)

		db.Conn.Where("Designation LIKE ?", "%"+job.Designation+"%").Find(&jobs)
		if len(resp) == 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return!!")
			return
		}

		json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	db.Conn.Model(&jobs).
		Joins("JOIN organizations ON jobs.org_id = organizations.org_id").
		Select("jobs.*, organizations.name").
		Scan(&resp)

	if len(resp) == 0 {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Nothing to return!!")
		return
	}

	json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
