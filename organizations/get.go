package organizations

import (
	"PostJson/db"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func GetOrganizations(w http.ResponseWriter, r *http.Request) {
	var organizations []Organizations
	var organization Organizations
	var add Organizations

	add.Org_ID = r.URL.Query().Get("id")
	if len(add.Org_ID) > 0 {
		// populating add for validation
		add.U_ID = "73c6ba9b-9325-4c68-bacb-52b6ce04e919"
		validate := validator.New()
		err := validate.Struct(add)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		err2 := db.Conn.QueryRow("SELECT org_id, name, about, website, u_id FROM organizations where org_id = ?", add.Org_ID).Scan(&organization.Org_ID, &organization.Name, &organization.About, &organization.Website, &organization.U_ID)
		if err2 != nil {
			w.WriteHeader(404)
			fmt.Fprintf(w, "Organization Not Found!!")
			return
		}
		json.Marshal(organization)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(organization)
		return
	}

	getAll, err := db.Conn.Query("select * from organizations")
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Record not found!!")
		return
	}
	for getAll.Next() {
		getAll.Scan(&organization.Org_ID, &organization.Name, &organization.About, &organization.Website, &organization.U_ID)
		organizations = append(organizations, organization)
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
