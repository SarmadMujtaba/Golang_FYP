package organizations

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func DeleteOrganizations(w http.ResponseWriter, r *http.Request) {
	var check structures.Organizations
	var organizations []structures.Organizations
	var members []structures.Memberships
	var jobs []structures.Jobs
	var reqSkills []structures.RequiredSkills
	var apps []structures.Applications

	wrongInput := true

	check.Org_ID = r.URL.Query().Get("id")
	if len(check.Org_ID) > 0 {
		// populating add for validation
		check.U_ID = check.Org_ID
		validate := validator.New()
		err := validate.Struct(check)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}
	}

	db.Conn.Find(&organizations)
	for _, usr := range organizations {
		if usr.Org_ID == check.Org_ID {
			db.Conn.Where("Org_id = ?", check.Org_ID).Find(&jobs)
			for _, row := range jobs {
				db.Conn.Where("ID = ?", row.ID).Delete(&reqSkills)
				db.Conn.Where("ID = ?", row.ID).Delete(&apps)
			}
			db.Conn.Where("Org_id = ?", check.Org_ID).Delete(&jobs)
			db.Conn.Where("Org_ID = ?", check.Org_ID).Delete(&members)
			db.Conn.Where("Org_ID = ?", check.Org_ID).Delete(&organizations)
			w.WriteHeader(200)
			fmt.Fprintf(w, "Organization deteled successfully!!")
			wrongInput = false
		}
	}

	if wrongInput {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Organization does not exist!!")
		return
	}
}
