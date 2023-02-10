package userprofile

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/go-playground/validator.v9"
)

// swagger:route PUT /profile Profile post-profile
//
// Add Profile
//
// You can add a user profile through this endpoint by filling in the details of the user to be added.
//
// responses:
//  200: Profile
//  201: Users
//  409: Error
//  400: Error

func Profile(w http.ResponseWriter, r *http.Request) {
	var profile structures.Profile
	var profiles []structures.Profile
	var user []structures.Users
	var skills []structures.Skills
	var experience []structures.Experience

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	profile.U_ID = dataToCompare["user_id"]
	profile.Education = dataToCompare["education"]
	profile.Phone = dataToCompare["phone"]

	validate := validator.New()
	err := validate.Struct(profile)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect Input")
		return
	}

	// validating json schema
	schemaLoader := gojsonschema.NewReferenceLoader("file:///app/schemas/ProfileSchema.json")
	documentLoader := gojsonschema.NewGoLoader(dataToCompare)

	res, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}
	if !res.Valid() {
		w.WriteHeader(400)
		for _, desc := range res.Errors() {
			fmt.Fprintln(w, desc.Description())
		}
		return
	}

	db.Conn.Where("ID = ?", profile.U_ID).Find(&user)
	if len(user) == 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "User does not exist!!")
		return
	}

	db.Conn.Where("U_ID = ?", profile.U_ID).Find(&profiles)
	if len(profiles) == 0 {
		result := db.Conn.Create(&profile)
		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Could not add profile!!")
			return
		}
		w.WriteHeader(201)
		fmt.Fprintf(w, "Profile added!!")
		return
	} else {
		// if profile already exists, delete its skills and exp, and update its profile (Put request)
		db.Conn.Where("U_ID = ?", profile.U_ID).Delete(&skills)
		db.Conn.Where("U_ID = ?", profile.U_ID).Delete(&experience)
		db.Conn.Model(&profile).Where("U_ID = ?", profile.U_ID).Update(&profile)
		w.WriteHeader(200)
		fmt.Fprintf(w, "Profile Updated!!")
		return
	}

}
