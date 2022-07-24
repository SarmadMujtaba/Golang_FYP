package jobs

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"
)

// swagger:route GET /category Jobs catParam
//
// Lists all jobs of a category
//
// This endpoint returns category specific jobs (1: Full-Time, 2: Part-Time, 3: Internship) if you pass that category ID as a query parameter
//
// responses:
//  200: Jobs
//  400: Error

func GetCategory(w http.ResponseWriter, r *http.Request) {
	var jobs []structures.Jobs
	var job structures.Jobs

	job.Cat_ID = r.URL.Query().Get("category_id")

	if len(job.Cat_ID) > 0 {
		if job.Cat_ID == "1" || job.Cat_ID == "2" || job.Cat_ID == "3" {
			db.Conn.Find(&jobs, "Cat_ID = ?", job.Cat_ID)
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
	}
	w.WriteHeader(400)
	fmt.Fprintf(w, "Incorrect Input!!")
}
