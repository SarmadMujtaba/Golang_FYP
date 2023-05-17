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
	UserID      string `json:"id"`
	Name        string `json:"name"`
	Education   string `json:"education"`
	Phone       string `json:"phone"`
	Experiences string `json:"experiences"`
	Skills      string `json:"skills"`
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	var profile structures.Profile
	var result Result

	profile.U_ID = strings.ReplaceAll(r.URL.Query().Get("user_id"), `"`, "")
	fmt.Println(profile.U_ID)
	if len(profile.U_ID) > 0 {
		err := db.Conn.Table("users").
			Joins("LEFT JOIN profiles ON users.id = profiles.u_id").
			Where("users.id = ?", profile.U_ID).
			Select("users.id, users.name, profiles.education, profiles.phone").
			Scan(&result).Error

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error retrieving profile: %v", err)
			return
		}

		if result.Phone == "" {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Nothing to return!!")
			return
		}

		// Fetch experiences
		var experiences []structures.Experience
		_ = db.Conn.Table("experiences").
			Where("u_id = ?", profile.U_ID).
			Select("experience").
			Find(&experiences).Error

		// Fetch skills
		var skills []structures.Skills
		_ = db.Conn.Table("skills").
			Where("u_id = ?", profile.U_ID).
			Select("skill").
			Find(&skills).Error

		// Populate experiences and skills in the result
		for _, exp := range experiences {
			result.Experiences += exp.Experience + ", "
		}

		for _, skill := range skills {
			result.Skills += skill.Skill + ", "
		}

		// Marshal the result into JSON
		jsonData, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error marshaling JSON: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	} else {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Missing or wrong Parameters!!")
	}
}
