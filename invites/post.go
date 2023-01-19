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
	"sync"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

var wg sync.WaitGroup

func SendEmail(uid string, orgid string, mail string, wr http.ResponseWriter) error {
	defer wg.Done()
	var owner structures.Users
	var org structures.Organizations

	db.Conn.Find(&owner, "ID = ?", uid)
	db.Conn.Find(&org, "Org_ID = ?", orgid)

	from := "191387@students.au.edu.pk"
	password := os.Getenv("EMAIL_PASSWORD")

	to := []string{
		mail,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("HTML_Templates/invite_template.html")

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
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		wr.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintln(wr, "Sending invitation link failed!!")
		return err
	} else {
		wr.WriteHeader(http.StatusOK)
		fmt.Fprintln(wr, "Invitation link sent on email!!")
		return nil
	}
}

func PostInvite(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var add structures.Invites
	var invites []structures.Invites
	var users []structures.Users
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
		// Sending invitation email to user through SendMail Goroutine

		uid := string(dataToCompare["owner"])
		orgid := string(dataToCompare["org_id"])
		mail := string(dataToCompare["target_email"])
		// errs channel will hold the error in case SendMail returns it
		// errs := make(chan error, 1)
		// go func() {
		// 	errs <- SendEmail(uid, orgid, mail, w)
		// }()
		wg.Add(1)
		go SendEmail(uid, orgid, mail, w)
		// if SendMail returned error (email sending failed), function will not signup
		// if err := <-errs; err != nil {
		// 	elapsed := time.Since(start)
		// 	fmt.Printf("Binomial took %s", elapsed)
		// 	return
		// }
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
				} else {
					w.WriteHeader(201)
					fmt.Fprintf(w, "Invite Added!!")
				}
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

	elapsed := time.Since(start)
	fmt.Printf("Binomial took %s", elapsed)
	wg.Wait()
}
