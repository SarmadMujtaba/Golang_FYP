package organizations

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func GetOrganizations(w http.ResponseWriter, r *http.Request) {
	var organizations []structures.Organizations
	var add structures.Organizations
	test := true

	add.Org_ID = r.URL.Query().Get("id")
	if len(add.Org_ID) > 0 {
		// // populating add for validation
		// add.U_ID = "73c6ba9b-9325-4c68-bacb-52b6ce04e919"
		// validate := validator.New()
		// err := validate.Struct(add)
		// if err != nil {
		// 	w.WriteHeader(400)
		// 	fmt.Fprintf(w, "Incorrect input!!")
		// 	return
		// }

		db.Conn.Find(&organizations) //QueryRow("SELECT org_id, name, about, website, u_id FROM organizations where org_id = ?", add.Org_ID).Scan(&organization.Org_ID, &organization.Name, &organization.About, &organization.Website, &organization.U_ID)
		for _, org := range organizations {
			if org.Org_ID == add.Org_ID {
				json.Marshal(org)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(org)
				test = false
				return
			}
		}

		result := db.Conn.Find(&organizations, "U_ID = ?", add.Org_ID)
		if result.Error == nil {
			json.Marshal(organizations)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(organizations)
			test = false
			return
		}

		if test == true {
			w.WriteHeader(400)
			fmt.Fprintf(w, "No organizations found!!")
			return
		}

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

	db.Conn.Find(&organizations) //Query("select * from organizations")

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
