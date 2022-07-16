package jobs

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	var check structures.Jobs
	var jobs []structures.Jobs
	var reqSkills []structures.RequiredSkills

	wrongInput := true

	check.ID = r.URL.Query().Get("id")
	if len(check.ID) > 0 {
		// populating add for validation
		check.Org_id = check.ID
		validate := validator.New()
		err := validate.Struct(check)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}
	}

	db.Conn.Find(&reqSkills)
	db.Conn.Find(&jobs)
	for _, usr := range jobs {
		if usr.ID == check.ID {
			db.Conn.Where("ID = ?", check.ID).Delete(&reqSkills)
			db.Conn.Where("ID = ?", check.ID).Delete(&jobs)
			w.WriteHeader(200)
			fmt.Fprintf(w, "Job deteled successfully!!")
			wrongInput = false
		}
	}

	if wrongInput {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Job does not exist!!")
		return
	}
}
