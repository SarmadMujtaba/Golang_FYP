package organizations

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

type Organizations struct {
	U_ID    string `json:"user_id" validate:"uuid"`
	Org_ID  string `json:"id" validate:"uuid"`
	Name    string `json:"name"`
	About   string `json:"about"`
	Website string `json:"website"`
}

func PostOrganizations(w http.ResponseWriter, r *http.Request) {
	var add Organizations

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	id := uuid.New()
	add.Org_ID = id.String()
	add.U_ID = dataToCompare["user_id"]
	add.Name = dataToCompare["name"]
	add.About = dataToCompare["about"]
	add.Website = dataToCompare["website"]

	// input validation
	validate := validator.New()
	err := validate.Struct(add)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	_, er := db.Conn.Query("insert into organizations values(?, ?, ?, ?, ?)", add.Org_ID, add.Name, add.About, add.Website, add.U_ID)
	if er != nil {
		fmt.Fprintln(w, "Could not enter record!!")
		return
	}

	id = uuid.New()
	_, err = db.Conn.Query("insert into membership values(?, ?, ?)", id.String(), add.U_ID, add.Org_ID)
	if err != nil {
		w.WriteHeader(400)
		db.Conn.Query("delete from organizations where org_id = ?", add.Org_ID)
		fmt.Fprintln(w, "Could not enter record!!")
		return
	}

	w.WriteHeader(201)
	fmt.Fprintf(w, "Organization Created!!")
	return
}
