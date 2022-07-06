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
)

func PostOrganizations(w http.ResponseWriter, r *http.Request) {
	var add structures.Organizations

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	id := uuid.New()
	add.Org_ID = id.String()
	add.U_ID = dataToCompare["user_id"]
	add.Name = dataToCompare["name"]
	add.About = dataToCompare["about"]
	add.Website = dataToCompare["website"]

	result := db.Conn.Create(&add) //Query("insert into organizations values(?, ?, ?, ?, ?)", add.Org_ID, add.Name, add.About, add.Website, add.U_ID)
	if result.Error != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not enter record!!")
		return
	}

	// id = uuid.New()
	// _, err = db.Conn.Query("insert into membership values(?, ?, ?)", id.String(), add.U_ID, add.Org_ID)
	// if err != nil {
	// 	w.WriteHeader(400)
	// 	db.Conn.Query("delete from organizations where org_id = ?", add.Org_ID)
	// 	fmt.Fprintln(w, "Could not enter record!!")
	// 	return
	// }

	w.WriteHeader(201)
	fmt.Fprintf(w, "Organization Created!!")
	return
}
