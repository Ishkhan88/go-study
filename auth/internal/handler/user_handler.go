package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Ishkhan88/go-study/auth/internal/service"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ProfileResponse struct {
	Message string `json:"message"`
	UserID  any    `json:"user_id"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "Данные пользователя"
// @Success 201 {object} RegisterResponse
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "cannot create user"
// @Router /api/users [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.Service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		http.Error(w, "cannot create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary Логин пользователя
// @Description Выполняет вход и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Логин пользователя"
// @Success 200 {object} LoginResponse
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "invalid email or password"
// @Router /api/auth/users [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.Service.LoginWithToken(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// Profile godoc
// @Summary Профиль пользователя
// @Description Возвращает user_id из JWT токена
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProfileResponse
// @Failure 401 {string} string "missing authorization header"
// @Router /api/users/profile [get]
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "profile доступен",
		"user_id": userID,
	})
}
