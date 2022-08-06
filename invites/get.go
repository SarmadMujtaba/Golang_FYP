package invites

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

func GetInvites(w http.ResponseWriter, r *http.Request) {
	var allInvites []structures.Invites
	var invites []structures.Invites
	var get structures.Invites

	get.Org_ID = r.URL.Query().Get("org_id")
	if len(get.Org_ID) > 0 {
		// populating add for validation
		get.U_ID = get.Org_ID
		validate := validator.New()
		err := validate.Struct(get)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		db.Conn.Find(&allInvites)
		for _, row := range allInvites {
			if row.Org_ID == get.Org_ID {
				invites = append(invites, row)
			}
		}

		if len(invites) > 0 {
			json.Marshal(invites)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invites)
			return
		} else {
			w.WriteHeader(400)
			fmt.Fprintf(w, "No invites found against given organizations!!")
			return
		}
	}

	get.Email = r.URL.Query().Get("target_email")
	if len(get.Email) > 0 {

		db.Conn.Find(&allInvites)
		for _, row := range allInvites {
			if row.Email == get.Email {
				invites = append(invites, row)
			}
		}

		if len(invites) > 0 {
			json.Marshal(invites)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invites)
			return
		} else {
			w.WriteHeader(400)
			fmt.Fprintf(w, "No invites found against given Email!!")
			return
		}
	}
	w.WriteHeader(400)
	fmt.Fprintf(w, "Missing or incorrect parameters!!")
}
