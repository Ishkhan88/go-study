package http

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

	expectedLogin := os.Getenv("LOGIN")
	expectedPassword := os.Getenv("PASSWORD")
	secret := os.Getenv("JWT_SECRET")

	if expectedLogin == "" || expectedPassword == "" || secret == "" {
		writeErr(w, 500, "auth env is not configured")
		return
	}

	if req.Login != expectedLogin || req.Password != expectedPassword {
		writeErr(w, 401, "invalid credentials")
		return
	}

	claims := jwt.MapClaims{
		"sub": req.Login,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		writeErr(w, 500, "token sign error")
		return
	}

	writeJSON(w, 200, LoginResponse{Token: signed})
}
