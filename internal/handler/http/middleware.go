package http

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		writeErr(w, 401, "missing bearer token")
		return false
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		writeErr(w, 401, "invalid authorization header")
		return false
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		writeErr(w, 500, "jwt secret is not configured")
		return false
	}

	tokenStr := parts[1]
	_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		writeErr(w, 401, "invalid token")
		return false
	}

	return true
}
