package users

import (
	"PostJson/db"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	var check Users
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

	var find string
	db.Conn.QueryRow("Select name from user_data where id = ?", check.ID).Scan(&find)
	if find == "" {
		w.WriteHeader(404)
		fmt.Fprintf(w, "User does not exist!!")
		return
	}
	var u_id string
	db.Conn.QueryRow("select id from membership where id = ?", check.ID).Scan(&u_id)
	db.Conn.Query("DELETE from membership where id = ?", check.ID)
	db.Conn.Query("DELETE from organizations where u_id = ?", u_id)
	db.Conn.Query("DELETE from user_data where id = ?", check.ID)
	w.WriteHeader(200)
	fmt.Fprintf(w, "Record deleted successfully!!")
	return
}
