package invites

import (
	"PostJson/db"
	"PostJson/structures"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"

	"gopkg.in/go-playground/validator.v9"
)

func PostInvite(w http.ResponseWriter, r *http.Request) {
	var add structures.Invites
	var invites []structures.Invites
	var users []structures.Users
	var owner structures.Users
	var org structures.Organizations
	var owners []structures.Organizations
	access := false
	ownerFound := false
	userFound := false

	dataFromWeb, _ := ioutil.ReadAll(r.Body)
	var dataToCompare map[string]string
	json.Unmarshal(dataFromWeb, &dataToCompare)

	add.U_ID = dataToCompare["owner"]
	add.Org_ID = dataToCompare["org_id"]
	add.Status = dataToCompare["status"]
	add.Email = dataToCompare["target_email"]

	validate := validator.New()
	err := validate.Struct(add)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Incorrect input!!")
		return
	}

	// validating json schema
	// schemaLoader := gojsonschema.NewReferenceLoader("file:///home/sarmad/Go_Practice/PostJson/schemas/OrgSchema.json")
	// documentLoader := gojsonschema.NewGoLoader(dataToCompare)

	// res, err := gojsonschema.Validate(schemaLoader, documentLoader)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// if !res.Valid() {
	// 	w.WriteHeader(400)
	// 	for _, desc := range res.Errors() {
	// 		fmt.Fprintln(w, desc.Description())
	// 	}
	// 	return
	// }
	db.Conn.Find(&invites)
	for _, row := range invites {
		if row.Email == add.Email && row.Org_ID == add.Org_ID {
			w.WriteHeader(409)
			fmt.Fprintln(w, "Invite Already Exist!!")
			return
		}
	}

	db.Conn.Find(&users)
	for _, row := range users {
		if row.Email == add.Email {
			userFound = true
		}
	}

	if !userFound {
		// Sending invitation email to user
		db.Conn.Find(&owner, "ID = ?", add.U_ID)
		db.Conn.Find(&org, "Org_ID = ?", add.Org_ID)

		from := "191387@students.au.edu.pk"
		password := os.Getenv("EMAIL_PASSWORD")

		to := []string{
			dataToCompare["target_email"],
		}

		smtpHost := "smtp.gmail.com"
		smtpPort := "587"

		auth := smtp.PlainAuth("", from, password, smtpHost)

		t, _ := template.ParseFiles("invite_template.html")

		var body bytes.Buffer

		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf("Subject: Invitation \n%s\n\n", mimeHeaders)))

		// Appending token to HTML file
		t.Execute(&body, struct {
			Owner    string
			Org_name string
		}{
			Owner:    owner.Name,
			Org_name: org.Name,
		})

		// Sending email.
		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Fprintln(w, "Sending invitation link failed!!")
			return
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Invitation link sent on email!!")
		}
	}

	db.Conn.Find(&owners)
	for _, onr := range owners {
		if onr.U_ID == add.U_ID {
			ownerFound = true
			if onr.Org_ID == add.Org_ID {
				access = true
				result := db.Conn.Create(&add)
				if result.Error != nil {
					w.WriteHeader(400)
					fmt.Fprintln(w, "Could not add invite!!")
					return
				}
				w.WriteHeader(201)
				fmt.Fprintf(w, "Invite Added!!")
			}
		}
	}

	if !access {
		w.WriteHeader(400)
		fmt.Fprintln(w, "This user does not have access to this organization!!")
		return
	}

	if !ownerFound {
		w.WriteHeader(400)
		fmt.Fprintln(w, "This user does not owns any organization!!")
		return
	}
}
