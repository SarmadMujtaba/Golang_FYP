package authentication

import (
	"PostJson/db"
	"PostJson/structures"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func Verify(w http.ResponseWriter, r *http.Request) {
	var user structures.Users

	token := r.URL.Query().Get("token")

	claims := &structures.Claims{}

	// Parsing and decoding token and storing in claims
	tkn, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		db.Conn.Model(&user).Where("email = ?", claims.Email).Delete(user)
		fmt.Fprintln(w, "Verification link Expired!!")
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Token Invalid!!")
		return
	}

	db.Conn.Model(&user).Where("email = ?", claims.Email).Update("IsVerified", true)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email Verified!!")
}
