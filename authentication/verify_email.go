package authentication

// using SendGrid's Go Library
// https://github.com/sendgrid/sendgrid-go

import (
	"PostJson/db"
	"PostJson/structures"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt"
)

func VerifyEmail(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var User structures.Users

		dataFromWeb, _ := ioutil.ReadAll(r.Body)
		var dataToCompare map[string]string
		json.Unmarshal(dataFromWeb, &dataToCompare)

		db.Conn.Find(&User, "email = ?", dataToCompare["email"])

		if User.Email == dataToCompare["email"] && User.IsVerified {
			w.WriteHeader(409)
			fmt.Fprintf(w, "Email ID already exist!!")
			return
		}

		// JWT token generation
		expirationTime := time.Now().Add(time.Minute * 10)

		claims := &structures.Claims{
			Email: dataToCompare["email"],
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

		//os.Setenv("Password", "DummyUniID")

		// Receiver email address.
		to := []string{
			dataToCompare["email"],
		}

		// smtp server configuration.
		smtpHost := "smtp.gmail.com"
		smtpPort := "587"

		// Authentication.
		auth := smtp.PlainAuth("", from, password, smtpHost)

		t, _ := template.ParseFiles("email_template.html")

		var body bytes.Buffer

		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf("Subject: Email Verification \n%s\n\n", mimeHeaders)))

		// Appending token to HTML file
		t.Execute(&body, struct {
			Token string
		}{
			Token: tokenString,
		})

		// Sending email.
		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Fprintln(w, "Sending verification link failed!!")
		} else {
			// Sending context of User data (JSON object) to next handler
			ctx := context.Background()
			ctx = context.WithValue(ctx, "object", dataToCompare)
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
