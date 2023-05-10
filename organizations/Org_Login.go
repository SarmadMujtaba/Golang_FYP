package organizations

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xeipuuv/gojsonschema"
)

// swagger:route GET /organizations Organization get-organizations
//
// Lists all / single organizations
//
// This endpoint returns all organizations if no query parameter is passed. However, it returns single organization if you pass its ID as a query parameter
//
// responses:
//  200: Organizations
//  404: Error
//  400: Error
func GetOrganizations(w http.ResponseWriter, r *http.Request) {
	var credentials structures.Organizations
	allOrgs := []structures.Organizations{}

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)
	fmt.Println(dataToCompare)

	credentials.Email = dataToCompare["email"]
	credentials.Pass = dataToCompare["pass"]

	// validating json schema
	schemaLoader := gojsonschema.NewReferenceLoader("file:///app/schemas/LoginSchema.json")
	documentLoader := gojsonschema.NewGoLoader(dataToCompare)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}
	if !result.Valid() {
		w.WriteHeader(400)
		fmt.Println("frontend error")
		fmt.Fprintf(w, "Json Object is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Fprintln(w, desc.Description())
		}
		return
	}

	db.Conn.Find(&allOrgs)
	for _, org := range allOrgs {
		if org.Email == credentials.Email && org.Pass == credentials.Pass {

			if !org.IsVerified {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Email not verified!!")
				return
			}

			// Generating Json Web Token for authenticating further requests

			// Token will be valid for one week
			// expirationTime := time.Now().Add((time.Hour * 24) * 7)

			// claims := &structures.Claims{
			// 	Email: credentials.Email,
			// 	StandardClaims: jwt.StandardClaims{
			// 		ExpiresAt: expirationTime.Unix(),
			// 	},
			// }

			// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			// tokenString, err := token.SignedString(jwtKey)

			// if err != nil {
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	return
			// }

			// http.SetCookie(w,
			// 	&http.Cookie{
			// 		Name:    "token",
			// 		Value:   tokenString,
			// 		Expires: expirationTime,
			// 		// HttpOnly will
			// 		HttpOnly: true,
			// 	})
			json.Marshal(org)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(org)
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "Incorrect Email or Password!!")
}
