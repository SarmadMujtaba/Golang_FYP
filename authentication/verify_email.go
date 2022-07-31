package authentication

// using SendGrid's Go Library
// https://github.com/sendgrid/sendgrid-go

import (
	"PostJson/structures"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func VerifyEmail(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		dataFromWeb, _ := ioutil.ReadAll(r.Body)
		var dataToCompare map[string]string
		json.Unmarshal(dataFromWeb, &dataToCompare)

		// Composing Email
		from := mail.NewEmail("J2E", "srmdmjtba@gmail.com")
		subject := "Email Verification"
		to := mail.NewEmail("User", dataToCompare["email"])
		plainTextContent := "Welcome to J2E, Please verify your email by clicking following link."

		// JWT token generation
		expirationTime := time.Now().Add(time.Minute * 2)

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

		htmlContent := `<a href="http://localhost:5020/verify?token=` + tokenString + `">Verify Email!</a>`
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

		// Passing SendGrid's API key as Parameter
		client := sendgrid.NewSendClient("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")
		_, err = client.Send(message)
		if err != nil {
			log.Println(err)
		} else {
			// Sending context of User data (JSON object) to next handler
			ctx := context.Background()
			ctx = context.WithValue(ctx, "object", dataToCompare)
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
