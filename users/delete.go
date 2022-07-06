package users

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	var check structures.Users
	var organizations []structures.Organizations
	var users []structures.Users
	test := true

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

	db.Conn.Find(&organizations)
	for _, orgs := range organizations {
		if orgs.U_ID == check.ID {
			db.Conn.Where("U_ID = ?", check.ID).Delete(&organizations)
			test = false
		}
	}

	db.Conn.Find(&users)

	for _, usr := range users {
		if usr.ID == check.ID {
			db.Conn.Where("ID = ?", check.ID).Delete(&users)
			w.WriteHeader(200)
			fmt.Fprintf(w, "User deteled successfully!!")
			test = false
		}
	}
	if test == true {
		w.WriteHeader(404)
		fmt.Fprintf(w, "User does not exist!!")
		return
	}

}
