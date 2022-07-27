package members

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/go-playground/validator.v9"
)

// swagger:route POST /members Memberships post-members
//
// Add membership
//
// You can add a member of an organization through this endpoint by filling in the details.
//
// responses:
//  201: Memberships
//  409: Error
//  400: Error

func PostMembers(w http.ResponseWriter, r *http.Request) {
	var members []structures.Memberships
	var member structures.Memberships
	duplicate := true

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	id := uuid.New()
	member.ID = id.String()
	member.U_ID = dataToCompare["user_id"]
	member.Org_ID = dataToCompare["org_id"]

	validate := validator.New()
	err := validate.Struct(member)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	// validating json schema
	schemaLoader := gojsonschema.NewReferenceLoader("file:///home/sarmad/Go_Practice/PostJson/schemas/MemberSchema.json")
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

	db.Conn.Find(&members)
	for _, row := range members {
		if row.U_ID == member.U_ID {
			if row.Org_ID == member.Org_ID {
				duplicate = false
			}
		}
	}

	if !duplicate {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Membership Already Exist!!")
		return
	}

	result := db.Conn.Create(&member)
	if result.Error != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not enter record!!")
		return
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, "Membership Added!!")
}
