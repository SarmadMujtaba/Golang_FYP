package jobs

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

// swagger:route GET /jobs?org_id Jobs jobParam
//
// Lists all / single Job
//
// This endpoint returns all Jobs if no query parameter is passed. However, it returns organization specific jobs if you pass that organization's ID as a query parameter
//
// responses:
//  200: Jobs
//  404: Error
//  400: Error

func GetJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []structures.Jobs
	var job structures.Jobs

	job.Org_id = r.URL.Query().Get("org_id")

	if len(job.Org_id) > 0 {
		// populating add for validation
		job.ID = job.Org_id
		validate := validator.New()
		err := validate.Struct(job)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}
		// Returns all jobs of an organization
		db.Conn.Find(&jobs, "Org_id = ?", job.Org_id)

		if len(jobs) == 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return!!")
			return
		}

		json.Marshal(jobs)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jobs)
		return
	}
	// Returns all jobs against a designation name
	job.Designation = r.URL.Query().Get("job_name")
	if len(job.Designation) > 0 {
		db.Conn.Where("Designation LIKE ?", "%"+job.Designation+"%").Find(&jobs)
		if len(jobs) == 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return!!")
			return
		}

		json.Marshal(jobs)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jobs)
		return
	}

	db.Conn.Find(&jobs)

	if len(jobs) == 0 {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Nothing to return!!")
		return
	}

	json.Marshal(jobs)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}
