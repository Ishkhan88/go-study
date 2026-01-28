package service

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		writeErr(w, http.StatusUnauthorized, "missing bearer token")
		return false
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		writeErr(w, http.StatusInternalServerError, "JWT_SECRET is empty")
		return false
	}

	tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		writeErr(w, http.StatusUnauthorized, "invalid token")
		return false
	}

	return true
}
