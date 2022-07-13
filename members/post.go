package members

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

func PostMembers(w http.ResponseWriter, r *http.Request) {
	// var user structures.Users
	// var org structures.Organizations
	var member structures.Memberships
	var repeatCheck []structures.Memberships

	// user.ID = r.URL.Query().Get("user_id")
	// org.Org_ID = r.URL.Query().Get("org_id")
	member.U_ID = r.URL.Query().Get("user_id")
	member.Org_ID = r.URL.Query().Get("org_id")

	if len(member.Org_ID) > 0 && len(member.U_ID) > 0 {

		validate := validator.New()
		err := validate.Struct(member)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		db.Conn.Where("U_ID = ?", member.U_ID).Find(&repeatCheck)
		for _, row := range repeatCheck {
			if row.Org_ID == member.Org_ID {
				w.WriteHeader(409)
				fmt.Fprintf(w, "Membership already exist!!")
				return
			}
		}

		id := uuid.New()
		member.ID = id.String()
		result := db.Conn.Create(&member)
		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Could not enter record!!")
			return
		}
		fmt.Fprintln(w, "Member Added!!")
		return
	}
}
