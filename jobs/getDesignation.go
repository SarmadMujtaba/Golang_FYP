package jobs

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"
)

// swagger:route GET /jobs/designations Jobs nameParam
//
// Lists Jobs of given designation
//
// This endpoint returns designation specific jobs if you pass that designation's name as a query parameter
//
// responses:
//  200: Jobs
//  404: Error
//  400: Error

func GetDesignations(w http.ResponseWriter, r *http.Request) {
	var jobs []structures.Jobs
	var job structures.Jobs

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
	} else {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Missing Parameters!!")
		return
	}
}
