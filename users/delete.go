package users

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

// swagger:route DELETE /users User delete-users
//
// Delete User
//
// This endpoint deletes a user if you pass its ID as a query parameter
//
// responses:
//  200: Error
//  400: Error

func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	var check structures.Users
	var users []structures.Users
	wrongInput := true

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

		db.Conn.Find(&users)

		for _, usr := range users {
			if usr.ID == check.ID {
				db.Conn.Where("ID = ?", check.ID).Delete(&users)
				w.WriteHeader(200)
				fmt.Fprintf(w, "User deteled successfully!!")
				wrongInput = false
			}
		}
		if wrongInput {
			w.WriteHeader(400)
			fmt.Fprintf(w, "User does not exist!!")
			return
		}

	} else {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Missing Parameters!!")
		return
	}
}
