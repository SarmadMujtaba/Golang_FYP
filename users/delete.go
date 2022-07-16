package users

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	var check structures.Users
	var organizations []structures.Organizations
	var users []structures.Users
	var members []structures.Memberships
	var jobs []structures.Jobs
	var reqSkills []structures.RequiredSkills

	wrongInput := true

	check.ID = r.URL.Query().Get("id")
	if len(check.ID) > 0 {
		// populating add for validation
		check.Name = "dummy"
		check.Email = "dummy@gmail.com"
		check.Pass = "dummy"
		validate := validator.New()
		err := validate.Struct(check)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}
	}

	db.Conn.Find(&organizations)
	for _, orgs := range organizations {
		if orgs.U_ID == check.ID {
			// Deleting User's organizations' jobs and their required skills
			db.Conn.Where("Org_id = ?", orgs.Org_ID).Find(&jobs)
			for _, row := range jobs {
				db.Conn.Where("ID = ?", row.ID).Delete(&reqSkills)
			}
			db.Conn.Where("Org_id = ?", orgs.Org_ID).Delete(&jobs)
			// Deleting User's memberships to all organizations
			db.Conn.Where("Org_ID = ?", orgs.Org_ID).Delete(&members)
			wrongInput = false
		}
	}

	db.Conn.Find(&members)
	for _, member := range members {
		if member.U_ID == check.ID {
			db.Conn.Where("U_ID = ?", check.ID).Delete(&members)
			wrongInput = false
		}
	}

	db.Conn.Find(&organizations)
	for _, orgs := range organizations {
		if orgs.U_ID == check.ID {
			db.Conn.Where("U_ID = ?", check.ID).Delete(&organizations)
			wrongInput = false
		}

		db.Conn.Find(&users)

		for _, usr := range users {
			if usr.ID == check.ID {
				db.Conn.Where("ID = ?", check.ID).Delete(&users)
				w.WriteHeader(200)
				fmt.Fprintf(w, "User deteled successfully!!")
				wrongInput = false
			}
		}
		if wrongInput {
			w.WriteHeader(400)
			fmt.Fprintf(w, "User does not exist!!")
			return
		}

	}
}
