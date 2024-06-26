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

// swagger:route POST /profile/experience Profile add-experience
//
// Add Experience
//
// You can add a user profile's experience through this endpoint by filling in the details of the user.
//
// responses:
//  201: Users
//  400: Error

func AddExperience(w http.ResponseWriter, r *http.Request) {

	// var exp structures.Experience
	var user []structures.Users

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]interface{}
	json.Unmarshal(dataFromWeb, &dataToCompare)

	UserID := dataToCompare["user_id"].(string)

	var experiences []string
	skillsRaw := dataToCompare["experiences"].([]interface{})
	for _, s := range skillsRaw {
		experiences = append(experiences, s.(string))
	}

	fmt.Println(UserID, experiences)

	for _, v := range experiences {
		exp := structures.Experience{
			U_ID:       strings.ReplaceAll(UserID, `"`, ""),
			Experience: v,
		}
		fmt.Println(exp.U_ID)
		fmt.Println(exp.Experience)

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

		db.Conn.Where("ID = ?", exp.U_ID).Find(&user)
		if len(user) == 0 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "User does not exist!!")
			return
		}

		result := db.Conn.Create(&exp)
		if result.Error != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Could not add Required skill!!")
			return
		}
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, "Experiences added!!")
}

// var exp structures.Experience
// var user []structures.Users

// dataFromWeb, _ := ioutil.ReadAll(r.Body)
// var dataToCompare map[string]string
// json.Unmarshal(dataFromWeb, &dataToCompare)

// exp.U_ID = strings.ReplaceAll(dataToCompare["user_id"], `"`, "")
// exp.Experience = dataToCompare["experience"]

// fmt.Println(exp.U_ID)
// fmt.Println(exp.Experience)

// // validate := validator.New()
// // err := validate.Struct(exp)
// // if err != nil {
// // 	w.WriteHeader(400)
// // 	fmt.Fprintf(w, "Incorrect Input")
// // 	return
// // }

// // validating json schema
// schemaLoader := gojsonschema.NewReferenceLoader("file:///app/schemas/ExperienceSchema.json")
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

// db.Conn.Where("ID = ?", exp.U_ID).Find(&user)
// if len(user) == 0 {
// 	w.WriteHeader(400)
// 	fmt.Fprintf(w, "User does not exist!!")
// 	return
// }

// result := db.Conn.Create(&exp)
// if result.Error != nil {
// 	w.WriteHeader(400)
// 	fmt.Fprintf(w, "Could not add experience!!")
// 	return
// }
// w.WriteHeader(201)
// fmt.Fprintf(w, "experience added!!")
// }
