package organizations

import (
	"PostJson/db"
	"PostJson/structures"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/go-playground/validator.v9"
)

// swagger:route POST /organizations Organization post-organizations
//
// Add organization
//
// You can add an organization through this endpoint by filling in the details of the organization to be added.
//
// responses:
//  201: Organizations
//  404: Error
//  400: Error
var jwtKey = []byte(os.Getenv("JWT_KEY"))

func PostOrganizations(w http.ResponseWriter, r *http.Request) {
	var add structures.Organizations
	var orgs []structures.Organizations
	// var member structures.Memberships
	// duplicate := true

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	id := uuid.New()
	add.Org_ID = id.String()
	add.Email = dataToCompare["email"]
	add.Pass = dataToCompare["pass"]
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

	// validating json schema
	schemaLoader := gojsonschema.NewReferenceLoader("file:///app/schemas/OrgSchema.json")
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

	db.Conn.Find(&orgs)
	for _, org := range orgs {
		if org.Email == add.Email {
			if !org.IsVerified {
				w.WriteHeader(401)
				fmt.Fprintf(w, "Email ID Unverified, please verify!!")
				return
			}
			w.WriteHeader(409)
			fmt.Fprintf(w, "Organization with this Email already exist!!")
			return
		}
	}

	result := db.Conn.Create(&add)
	if result.Error != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not enter record!!")
		return
	} else {
		go SendEmail(dataToCompare["email"])
		w.WriteHeader(201)
		fmt.Fprintln(w, "Organization inserted, please visit your email for verification!!")
	}
}

func SendEmail(email string) {
	expirationTime := time.Now().Add(time.Minute * 30)

	fmt.Println(email)

	claims := &structures.Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Signing Token via key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	// Sender data.
	from := "191387@students.au.edu.pk"
	password := os.Getenv("EMAIL_PASSWORD")

	receiver := strings.ReplaceAll(email, " ", "")
	// Receiver email address.
	to := []string{
		receiver,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("HTML_Templates/org_verify.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: J2E Email Verification \n%s\n\n", mimeHeaders)))

	// Appending token to HTML file
	t.Execute(&body, struct {
		Token string
	}{
		Token: tokenString,
	})

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println("Sending Email Failed!!")
	}
}
