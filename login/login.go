package login

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

var jwtKey = []byte("secret_key")

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials structures.Users
	allUsers := []structures.Users{}

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	credentials.Email = dataToCompare["email"]
	credentials.Pass = dataToCompare["pass"]

	// validating json schema
	// schemaLoader := gojsonschema.NewReferenceLoader("file:///home/sarmad/Go_Practice/PostJson/structures/UserSchema.json")
	// documentLoader := gojsonschema.NewGoLoader(dataToCompare)

	// result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// if !result.Valid() {
	// 	w.WriteHeader(400)
	// 	fmt.Fprintf(w, "Json Object is not valid. see errors :\n")
	// 	for _, desc := range result.Errors() {
	// 		fmt.Fprintln(w, desc.Description())
	// 	}
	// 	return
	// }

	db.Conn.Find(&allUsers)
	for _, usr := range allUsers {
		if usr.Email == credentials.Email && usr.Pass == credentials.Pass {

			// Generating Json Web Token for authenticating further requests

			// Token will be valid for one week
			expirationTime := time.Now().Add((time.Hour * 24) * 7)

			claims := &structures.Claims{
				Email: credentials.Email,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w,
				&http.Cookie{
					Name:    "token",
					Value:   tokenString,
					Expires: expirationTime,
					// HttpOnly will
					HttpOnly: true,
				})

			w.WriteHeader(200)
			fmt.Fprintf(w, "Login Successful!!")
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "Incorrect Email or Password!!")
}
