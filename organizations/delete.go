package organizations

import (
	"PostJson/db"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func DeleteOrganizations(w http.ResponseWriter, r *http.Request) {
	var check Organizations
	check.Org_ID = r.URL.Query().Get("id")
	if len(check.Org_ID) > 0 {
		// populating add for validation
		check.U_ID = "73c6ba9b-9325-4c68-bacb-52b6ce04e919"
		validate := validator.New()
		err := validate.Struct(check)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}
	}

	var find string
	db.Conn.QueryRow("Select name from organizations where org_id = ?", check.Org_ID).Scan(&find)
	if find == "" {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Organization does not exist!!")
		return
	}
	db.Conn.Query("DELETE from membership where org_id = ?", check.Org_ID)
	db.Conn.Query("DELETE from organizations where org_id = ?", check.Org_ID)
	w.WriteHeader(200)
	fmt.Fprintf(w, "Record deleted successfully!!")
	return
}

func OrganizationsOfUser(w http.ResponseWriter, r *http.Request) {
	var users []Organizations
	var user Organizations
	var check Organizations

	check.U_ID = r.URL.Query().Get("id")
	if len(check.U_ID) > 0 {
		// populating add for validation
		check.Org_ID = "73c6ba9b-9325-4c68-bacb-52b6ce04e919"
		validate := validator.New()
		err := validate.Struct(check)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}
	}

	data, err := db.Conn.Query("select * from organizations where u_id = ?", check.U_ID)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "No organizations found against this user ID!!")
	}
	for data.Next() {
		data.Scan(&user.Org_ID, &user.Name, &user.About, &user.Website, &user.U_ID)
		users = append(users, user)
	}

	if len(users) == 0 {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Nothing to return!!")
		return
	}

	json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	return
}
