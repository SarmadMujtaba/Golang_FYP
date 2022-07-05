package users

import (
	"PostJson/db"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []Users
	var user Users
	var add Users

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

		err2 := db.Conn.QueryRow("SELECT * FROM user_data where id = ?", add.ID).Scan(&user.ID, &user.Email, &user.Name, &user.Pass)
		if err2 != nil {
			w.WriteHeader(404)
			fmt.Fprintf(w, "User Not Found!!")
			return
		}
		json.Marshal(user)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return
	}

	getAll, err := db.Conn.Query("select * from user_data")
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Failed to Fetch!!")
		return
	}
	for getAll.Next() {
		getAll.Scan(&user.ID, &user.Email, &user.Name, &user.Pass)
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
