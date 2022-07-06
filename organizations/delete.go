package organizations

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func DeleteOrganizations(w http.ResponseWriter, r *http.Request) {
	var check structures.Organizations
	var organizations []structures.Organizations
	test := true

	check.Org_ID = r.URL.Query().Get("id")
	// if len(check.Org_ID) > 0 {
	// 	// populating add for validation
	// 	check.U_ID = "73c6ba9b-9325-4c68-bacb-52b6ce04e919"
	// 	validate := validator.New()
	// 	err := validate.Struct(check)
	// 	if err != nil {
	// 		w.WriteHeader(400)
	// 		fmt.Fprintf(w, "Incorrect input!!")
	// 		return
	// 	}
	// }

	db.Conn.Find(&organizations)
	for _, usr := range organizations {
		if usr.Org_ID == check.Org_ID {
			db.Conn.Where("Org_ID = ?", check.Org_ID).Delete(&organizations)
			w.WriteHeader(200)
			fmt.Fprintf(w, "Organization deteled successfully!!")
			test = false
		}
	}

	if test == true {
		w.WriteHeader(404)
		fmt.Fprintf(w, "User does not exist!!")
		return
	}

	// var find string
	// db.Conn.QueryRow("Select name from organizations where org_id = ?", check.Org_ID).Scan(&find)
	// if find == "" {
	// 	w.WriteHeader(404)
	// 	fmt.Fprintf(w, "Organization does not exist!!")
	// 	return
	// }
	// db.Conn.Query("DELETE from membership where org_id = ?", check.Org_ID)
	// db.Conn.Query("DELETE from organizations where org_id = ?", check.Org_ID)
	// w.WriteHeader(200)
	// fmt.Fprintf(w, "Record deleted successfully!!")
	// return
}
