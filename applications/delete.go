package applications

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"
	"strings"
)

func DeleteApplications(w http.ResponseWriter, r *http.Request) {
	var app structures.Applications
	var apps []structures.Applications
	wrongInput := true

	app.U_ID = strings.ReplaceAll(r.URL.Query().Get("user_id"), `"`, "")
	app.Job_ID = strings.ReplaceAll(r.URL.Query().Get("job_id"), `"`, "")
	if len(app.U_ID) > 0 {
		// populating add for validation
		// app.Job_ID = app.U_ID
		// validate := validator.New()
		// err := validate.Struct(app)
		// if err != nil {
		// 	w.WriteHeader(400)
		// 	fmt.Fprintf(w, "Incorrect input!!")
		// 	return
		// }
		fmt.Println(app.U_ID)
		fmt.Println(app.Job_ID)

		db.Conn.Model(&apps).Where("U_ID = ? AND Job_ID = ?", app.U_ID, app.Job_ID).Delete(&apps)

		w.WriteHeader(200)
		fmt.Fprintf(w, "Application Deleted!!")
		wrongInput = false
		return

	}

	if wrongInput {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not delete application!!")
		return
	}

}
