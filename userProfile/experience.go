package userprofile

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

// swagger:route POST /profile/experience Profile post-experience
//
// Add Experience
//
// You can add a user profile's experience through this endpoint by filling in the details of the user.
//
// responses:
//  201: Users
//  400: Error

func AddExperience(w http.ResponseWriter, r *http.Request) {
	var exp structures.Experience
	var user []structures.Users

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	exp.U_ID = dataToCompare["user_id"]
	exp.Experience = dataToCompare["experience"]

	validate := validator.New()
	err := validate.Struct(exp)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect Input")
		return
	}

	db.Conn.Where("ID = ?", exp.U_ID).Find(&user)
	if len(user) == 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "User does not exist!!")
		return
	}

	result := db.Conn.Create(&exp)
	if result.Error != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Could not add experience!!")
		return
	}
	w.WriteHeader(201)
	fmt.Fprintf(w, "experience added!!")
}
