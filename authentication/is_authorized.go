package authentication

import (
	"PostJson/structures"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(401)
				fmt.Fprintf(w, "No cookie found!!")
				return
			}
			fmt.Fprintf(w, "Invalid or Missing token!!")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value

		claims := &structures.Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims,
			func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				fmt.Fprintf(w, "Request UnAuthorized!!")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			fmt.Fprintf(w, "Request UnAuthorized!!")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// refresh token after every 30 minutes
		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < (((time.Hour * 24) * 7) - (time.Minute * 30)) {
			expirationTime := time.Now().Add(((time.Hour * 24) * 7))

			claims.ExpiresAt = expirationTime.Unix()

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w,
				&http.Cookie{
					Name:     "token",
					Value:    tokenString,
					HttpOnly: true,
					Expires:  expirationTime,
				})
		}

		handler.ServeHTTP(w, r)
	}
}
