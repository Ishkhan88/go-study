package service

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler godoc
// @Summary Login
// @Description Login and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, 405, "only POST")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "bad json")
		return
	}

	// берем логин/пароль из .env
	adminLogin := os.Getenv("ADMIN_LOGIN")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		writeErr(w, 500, "JWT_SECRET is empty")
		return
	}

	// сравниваем
	if req.Login != adminLogin || req.Password != adminPassword {
		writeErr(w, 401, "invalid login or password")
		return
	}

	// создаём токен на 1 час
	claims := jwt.MapClaims{
		"sub": req.Login,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		writeErr(w, 500, "failed to sign token")
		return
	}

	writeJSON(w, 200, LoginResponse{Token: signedToken})
}
