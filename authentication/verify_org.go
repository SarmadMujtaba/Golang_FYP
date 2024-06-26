package authentication

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func VerifyOrg(w http.ResponseWriter, r *http.Request) {
	var org structures.Organizations

	token := r.URL.Query().Get("token")

	claims := &structures.Claims{}

	// Parsing and decoding token and storing in claims
	tkn, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		db.Conn.Find(&org, "email = ?", claims.Email)
		if !org.IsVerified {
			db.Conn.Model(&org).Where("email = ?", claims.Email).Delete(org)
		}

		http.ServeFile(w, r, "HTML_Templates/failed.html")
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Token Invalid!!")
		return
	}

	db.Conn.Model(&org).Where("email = ?", claims.Email).Update("IsVerified", true)
	http.ServeFile(w, r, "HTML_Templates/success.html")
}
