package organizations

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

func PostOrganizations(w http.ResponseWriter, r *http.Request) {
	var add structures.Organizations
	var users []structures.Users
	var member structures.Memberships
	duplicate := true

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	id := uuid.New()
	add.Org_ID = id.String()
	add.U_ID = dataToCompare["user_id"]
	add.Name = dataToCompare["name"]
	add.About = dataToCompare["about"]
	add.Website = dataToCompare["website"]

	validate := validator.New()
	err := validate.Struct(add)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	db.Conn.Find(&users)
	for _, usr := range users {
		if usr.ID == add.U_ID {
			duplicate = false
		}
	}

	if duplicate {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not enter record!!")
		return
	}

	result := db.Conn.Create(&add)
	if result.Error != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not enter record!!")
		return
	}

	id = uuid.New()
	member.ID = id.String()
	member.Org_ID = add.Org_ID
	member.U_ID = add.U_ID
	result = db.Conn.Create(&member)
	if result.Error != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not enter record!!")
		return
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, "Organization Created!!")
}
