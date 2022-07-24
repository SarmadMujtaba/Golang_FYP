package members

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

// swagger:route GET /members Memberships memberParam
//
// Lists Members
//
// This endpoint returns all members of an organizations based on organization's ID passed as query parameter.
//
// responses:
//  200: Memberships
//  400: Error

func GetMembers(w http.ResponseWriter, r *http.Request) {
	var member []structures.Memberships
	var add structures.Memberships
	var user []structures.Users
	wrongInput := true

	add.Org_ID = r.URL.Query().Get("org_id")
	if len(add.Org_ID) > 0 {
		// populating add for validation
		add.U_ID = add.Org_ID
		validate := validator.New()
		err := validate.Struct(add)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		isEmpty := db.Conn.Find(&member)
		if isEmpty.Value == nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "No organizations found!!")
			return
		}
		for _, member := range member {
			if member.Org_ID == add.Org_ID {
				wrongInput = false
				result := db.Conn.Find(&user, "ID = ?", member.U_ID)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result.Value)
			}
		}
		if wrongInput {
			w.WriteHeader(400)
			fmt.Fprintf(w, "No organizations found!!")
			return
		}
	} else {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Missing Parameters!!")
		return
	}
}
