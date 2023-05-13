package userprofile

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Result struct {
	UserID     string                  `json:"id"`
	Name       string                  `json:"name"`
	Education  string                  `json:"education"`
	Phone      string                  `json:"phone"`
	Experience []structures.Experience `json:"experiences"`
	Skills     []structures.Skills     `json:"skills"`
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	var profile structures.Profile
	// var exps []structures.Experience
	// var skills []structures.Skills
	var result Result

	profile.U_ID = strings.ReplaceAll(r.URL.Query().Get("user_id"), `"`, "")
	fmt.Println(profile.U_ID)
	if len(profile.U_ID) > 0 {
		// populating add for validation
		profile.Phone = "123"
		// validate := validator.New()
		// err := validate.Struct(profile)
		// if err != nil {
		// 	w.WriteHeader(400)
		// 	fmt.Fprintf(w, "Incorrect input!!")
		// 	return
		// }

		db.Conn.Table("users").
			Joins("LEFT JOIN profiles ON users.id = profiles.u_id").
			Joins("LEFT JOIN experiences ON profiles.u_id = experiences.u_id").
			Joins("LEFT JOIN skills ON profiles.u_id = skills.u_id").
			Where("users.id = ?", profile.U_ID).
			Select("users.id, users.name, profiles.education, profiles.phone, experiences.experience, skills.skill").
			Scan(&result)

		if result.Phone == "" {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return!!")
			return
		}

		// db.Conn.Find(&exps, "U_ID = ?", profile.U_ID)
		// db.Conn.Find(&skills, "U_ID = ?", profile.U_ID)

		// profile.Experience = exps
		// profile.Skills = skills

		json.Marshal(result)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	} else {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Missing or wrong Parameters!!")
	}

}
