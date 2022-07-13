package members

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

func GetMembers(w http.ResponseWriter, r *http.Request) {
	var member []structures.Memberships
	var add structures.Memberships
	var user []structures.Users
	wrongInput := true

	add.Org_ID = r.URL.Query().Get("id")
	if len(add.Org_ID) > 0 {
		// populating add for validation
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
				return
			}
		}
		if wrongInput == true {
			w.WriteHeader(400)
			fmt.Fprintf(w, "No organizations found!!")
			return
		}
	}
}
