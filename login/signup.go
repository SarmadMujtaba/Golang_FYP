package login

import (
	"PostJson/db"
	"PostJson/structures"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/go-playground/validator.v9"
)

// swagger:route POST /users User post-users
//
// Add User
//
// You can add a user through this endpoint by filling in the details of the user to be added.
//
// responses:
//  201: Users
//  409: Error
//  400: Error

func Signup(w http.ResponseWriter, r *http.Request) {
	var add structures.Users
	allUsers := []structures.Users{}

	// Receiving context of user data from middleware
	dataFromWeb := r.Context().Value("object").(map[string]string)

	id := uuid.New()
	add.ID = id.String()
	add.IsVerified = false
	add.Email = dataFromWeb["email"]
	add.Name = dataFromWeb["name"]
	add.Pass = dataFromWeb["pass"]

	// input validation
	validate := validator.New()
	err := validate.Struct(add)
	if err != nil {
		w.WriteHeader(400)
		fmt.Println("msla-1")
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	// validating json schema
	// without docker, replace the address with local address i.e "file:///home/sarmad/Go_Practice/PostJson/schemas/SignupSchema.json"

	schemaLoader := gojsonschema.NewReferenceLoader("file:///app/schemas/SignupSchema.json")
	documentLoader := gojsonschema.NewGoLoader(dataFromWeb)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}
	if !result.Valid() {
		w.WriteHeader(400)
		fmt.Println("msla-2")
		fmt.Fprintf(w, "Json Object is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Fprintln(w, desc.Description())
		}
		return
	}

	db.Conn.Find(&allUsers)
	for _, usr := range allUsers {
		if usr.Email == add.Email {
			if !usr.IsVerified {
				w.WriteHeader(401)
				fmt.Fprintf(w, "Email ID Unverified, please verify!!")
				return
			}
			w.WriteHeader(409)
			fmt.Fprintf(w, "Email ID already exist!!")
			return
		}
	}

	db.Conn.Create(&add)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(add)

	// fmt.Fprintf(w, "User inserted, please visit your email for verification!!")
}
