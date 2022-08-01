package login

import (
	"PostJson/db"
	"PostJson/structures"
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
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	// validating json schema
	schemaLoader := gojsonschema.NewReferenceLoader("file:///home/sarmad/Go_Practice/PostJson/schemas/UserSchema.json")
	documentLoader := gojsonschema.NewGoLoader(dataFromWeb)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}
	if !result.Valid() {
		w.WriteHeader(400)
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
	w.WriteHeader(201)
	fmt.Fprintf(w, "User inserted, please visit your email for verification!!")
}
