package userprofile

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// swagger:route POST /profile/skills Profile add-skills
//
// Add Skill
//
// You can add a user profile's Skills through this endpoint by filling in the details of the user.
//
// responses:
//  201: Users
//  400: Error

func AddSkill(w http.ResponseWriter, r *http.Request) {
	var user []structures.Users

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]interface{}
	json.Unmarshal(dataFromWeb, &dataToCompare)

	userID := dataToCompare["user_id"].(string)

	var skills []string
	skillsRaw := dataToCompare["skills"].([]interface{})
	for _, s := range skillsRaw {
		skills = append(skills, s.(string))
	}

	fmt.Println(userID, skills)

	for _, v := range skills {
		skill := structures.Skills{
			U_ID:  strings.ReplaceAll(userID, `"`, ""),
			Skill: v,
		}
		fmt.Println(skill.U_ID)
		fmt.Println(skill.Skill)

		// validate := validator.New()
		// err := validate.Struct(skill)
		// if err != nil {
		//  w.WriteHeader(400)
		//  fmt.Fprintf(w, "Incorrect Input")
		//  return
		// }

		// validating json schema
		// schemaLoader := gojsonschema.NewReferenceLoader("file:///app/schemas/ReqSkillSchema.json")
		// documentLoader := gojsonschema.NewGoLoader(dataToCompare)

		// res, err := gojsonschema.Validate(schemaLoader, documentLoader)
		// if err != nil {
		// 	panic(err.Error())
		// }
		// if !res.Valid() {
		// 	w.WriteHeader(400)
		// 	for _, desc := range res.Errors() {
		// 		fmt.Fprintln(w, desc.Description())
		// 	}
		// 	return
		// }

		db.Conn.Where("ID = ?", skill.U_ID).Find(&user)
		if len(user) == 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "User does not exist!!")
			return
		}

		result := db.Conn.Create(&skill)
		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Could not add Required skill!!")
			return
		}
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, "Skills added!!")
}

// package userprofile

// import (
// 	"PostJson/db"
// 	"PostJson/structures"
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"

// 	"github.com/xeipuuv/gojsonschema"
// 	"gopkg.in/go-playground/validator.v9"
// )

// // swagger:route POST /profile/skills Profile add-skills
// //
// // Add Skill
// //
// // You can add a user profile's Skills through this endpoint by filling in the details of the user.
// //
// // responses:
// //  201: Users
// //  400: Error

// func AddSkill(w http.ResponseWriter, r *http.Request) {
// 	var skill structures.Skills
// 	var user []structures.Users
// 	var skills []structures.Skills
// 	// var Skills []string

// 	dataFromWeb, _ := ioutil.ReadAll(r.Body)
// 	var dataToCompare map[string]string
// 	json.Unmarshal(dataFromWeb, &dataToCompare)

// 	skill.U_ID = strings.ReplaceAll(dataToCompare["user_id"], `"`, "")
// 	skill.Skill = dataToCompare["skill"]

// 	validate := validator.New()
// 	err := validate.Struct(skill)
// 	if err != nil {
// 		w.WriteHeader(400)
// 		fmt.Fprintf(w, "Incorrect Input")
// 		return
// 	}

// 	// validating json schema
// 	schemaLoader := gojsonschema.NewReferenceLoader("file:///app/schemas/SkillSchema.json")
// 	documentLoader := gojsonschema.NewGoLoader(dataToCompare)

// 	res, err := gojsonschema.Validate(schemaLoader, documentLoader)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	if !res.Valid() {
// 		w.WriteHeader(400)
// 		for _, desc := range res.Errors() {
// 			fmt.Fprintln(w, desc.Description())
// 		}
// 		return
// 	}

// 	db.Conn.Where("ID = ?", skill.U_ID).Find(&user)
// 	if len(user) == 0 {
// 		w.WriteHeader(400)
// 		fmt.Fprintf(w, "User does not exist!!")
// 		return
// 	}

// 	result := db.Conn.Create(&skill)
// 	if result.Error != nil {
// 		w.WriteHeader(400)
// 		fmt.Fprintf(w, "Could not add skill!!")
// 		return
// 	}

// 	// sending skills to python_db
// 	result2 := db.Conn.Model(&skills).Select("skill").Find(&skills)
// 	if result2.Error != nil {
// 		w.WriteHeader(400)
// 		fmt.Fprintf(w, "Nothing to return")
// 		return
// 	}

// 	data, _ := json.Marshal(dataToCompare["skill"])

// 	// change url with python's url later. It is Path parameter after url
// 	posturl := "http://host.docker.internal:8000/skills/" + skill.U_ID

// 	// concurently sending request to python
// 	go SendRequest(posturl, data)

// 	w.WriteHeader(201)
// 	fmt.Fprintf(w, "Skill added!!")
// }

// func SendRequest(url string, data []byte) {
// 	r, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	r.Header.Add("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(resp.StatusCode)
// }
