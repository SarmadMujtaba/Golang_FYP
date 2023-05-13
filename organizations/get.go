package organizations

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func GetOrg(w http.ResponseWriter, r *http.Request) {

	var add structures.Organizations
	var orgs []structures.Organizations
	// var member structures.Memberships
	duplicate := true

	add.Org_ID = strings.ReplaceAll(r.URL.Query().Get("id"), `"`, "")

	db.Conn.Find(&orgs)
	for _, org := range orgs {
		if org.Org_ID == add.Org_ID {
			json.Marshal(org)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(org)
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
