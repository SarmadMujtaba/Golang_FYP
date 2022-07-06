package users

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []structures.Users
	var add structures.Users
	var test bool = true

	add.ID = r.URL.Query().Get("id")
	if len(add.ID) > 0 {
		// populating add for validation
		add.Name = "dummy"
		add.Pass = "dummy"
		add.Email = "dummy@gmail.com"
		validate := validator.New()
		err := validate.Struct(add)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		db.Conn.Find(&users)
		for _, usr := range users {
			if usr.ID == add.ID {
				json.Marshal(usr)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(usr)
				test = false
				return
			}
		}
		if test == true {
			w.WriteHeader(404)
			fmt.Fprintf(w, "User Not Found!!")
			return
		}
	}

	db.Conn.Find(&users)

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
