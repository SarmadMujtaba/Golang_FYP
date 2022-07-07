package organizations

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func GetOrganizations(w http.ResponseWriter, r *http.Request) {
	var organizations []structures.Organizations
	var add structures.Organizations
	test := true

	add.Org_ID = r.URL.Query().Get("id")
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

		db.Conn.Find(&organizations)
		for _, org := range organizations {
			if org.Org_ID == add.Org_ID {
				test = false
				json.Marshal(org)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(org)
				return
			}
		}

		result := db.Conn.Find(&organizations, "U_ID = ?", add.Org_ID)
		if result.Value != nil {
			test = false

			// return
		}

		if test == true {
			w.WriteHeader(400)
			fmt.Fprintf(w, "No organizations found!!")
			return
		}

		if len(organizations) == 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return!!")
			return
		}
		json.Marshal(organizations)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(organizations)
		return
	}

	db.Conn.Find(&organizations)

	if len(organizations) == 0 {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Nothing to return!!")
		return
	}

	json.Marshal(organizations)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(organizations)
	return
}
