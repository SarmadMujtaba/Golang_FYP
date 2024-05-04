package organizations

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func UpdateOrg(w http.ResponseWriter, r *http.Request) {

	var add structures.Organizations
	var orgs []structures.Organizations
	// var member structures.Memberships
	duplicate := true

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	add.Org_ID = strings.ReplaceAll(r.URL.Query().Get("id"), `"`, "")
	add.Email = dataToCompare["email"]
	add.Pass = dataToCompare["pass"]
	add.Name = dataToCompare["name"]
	add.About = dataToCompare["about"]
	add.Website = dataToCompare["website"]

	fmt.Println(add)

	db.Conn.Find(&orgs)
	for _, org := range orgs {
		if org.Org_ID == add.Org_ID {
			db.Conn.Model(&add).Where("Org_ID = ?", org.Org_ID).Update(add)
			w.WriteHeader(200)
			fmt.Fprintf(w, "Profile Updated successfully!!")
			duplicate = false
			return
		}
	}
	if duplicate {
		w.WriteHeader(404)
		fmt.Fprintf(w, "No organization with this ID!!")
		return
	}
}
