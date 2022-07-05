package users

import (
	"PostJson/db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

type Users struct {
	ID    string `json:"id" validate:"uuid"`
	Name  string `json:"name" validate:"alpha"`
	Email string `json:"email" validate:"email"`
	Pass  string `json:"pass" validate:"alphanum"`
}

func PostUsers(w http.ResponseWriter, r *http.Request) {
	var add Users

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	id := uuid.New()
	add.ID = id.String()
	add.Email = dataToCompare["email"]
	add.Name = dataToCompare["name"]
	add.Pass = dataToCompare["pass"]

	// input validation
	validate := validator.New()
	err := validate.Struct(add)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	var duplicate string
	db.Conn.QueryRow("Select name from user_data where email = ?", add.Email).Scan(&duplicate)
	if duplicate != "" {
		w.WriteHeader(409)
		fmt.Fprintf(w, "Email ID already exist!!")
		return
	}

	_, err = db.Conn.Query("insert into user_data values(?, ?, ?, ?)", add.ID, add.Email, add.Name, add.Pass)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(201)
	fmt.Fprintf(w, "User inserted!!")
	return
}
