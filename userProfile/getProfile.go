package userprofile

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

// swagger:route GET /profile?user_id Profile userProfile
//
// Get User's Profile
//
// This endpoint returns a user's Profile if you pass its ID as a query parameter
//
// responses:
//  200: Users
//  400: Error

func GetProfile(w http.ResponseWriter, r *http.Request) {
	var profile structures.Profile
	var exps []structures.Experience
	var skills []structures.Skills

	profile.U_ID = r.URL.Query().Get("user_id")

	if len(profile.U_ID) > 0 {
		// populating add for validation
		profile.Phone = "123"
		validate := validator.New()
		err := validate.Struct(profile)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Incorrect input!!")
			return
		}

		db.Conn.Model(&profile).Where("U_ID = ?", profile.U_ID).Find(&profile)

		if profile.Education == "" {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return!!")
			return
		}

		db.Conn.Find(&exps, "U_ID = ?", profile.U_ID)
		db.Conn.Find(&skills, "U_ID = ?", profile.U_ID)

		profile.Experience = exps
		profile.Skills = skills

		json.Marshal(profile)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	} else {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Missing or wrong Parameters!!")
	}

}
