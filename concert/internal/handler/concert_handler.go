package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ishkhan88/go-study/concert/internal/model"
	"github.com/Ishkhan88/go-study/concert/internal/service"
)

type ConcertHandler struct {
	Service *service.ConcertService
}

func NewConcertHandler(s *service.ConcertService) *ConcertHandler {
	return &ConcertHandler{Service: s}
}

// GetConcerts godoc
// @Summary Список концертов
// @Description Возвращает список всех концертов
// @Tags concerts
// @Produce json
// @Success 200 {array} model.Concert
// @Failure 500 {string} string "cannot get concerts"
// @Router /concerts [get]
func (h *ConcertHandler) GetConcerts(w http.ResponseWriter, r *http.Request) {
	concerts, err := h.Service.GetAllConcerts()
	if err != nil {
		http.Error(w, "cannot get concerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(concerts)
}

// GetConcertByID godoc
// @Summary Получить концерт по id
// @Description Возвращает один концерт
// @Tags concerts
// @Produce json
// @Param id path int true "Concert ID"
// @Success 200 {object} model.Concert
// @Failure 400 {string} string "invalid concert id"
// @Failure 404 {string} string "concert not found"
// @Router /concerts/{id} [get]
func (h *ConcertHandler) GetConcertByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/concerts/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid concert id", http.StatusBadRequest)
		return
	}

	concert, err := h.Service.GetConcertByID(id)
	if err != nil {
		http.Error(w, "concert not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(concert)
}

// CreateConcert godoc
// @Summary Создать концерт
// @Description Создаёт новый концерт
// @Tags concerts
// @Accept json
// @Produce json
// @Param input body model.Concert true "Concert data"
// @Success 201 {object} model.Concert
// @Failure 400 {string} string "invalid request body"
// @Failure 500 {string} string "cannot create concert"
// @Router /concerts [post]
func (h *ConcertHandler) CreateConcert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var concert model.Concert

	err := json.NewDecoder(r.Body).Decode(&concert)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if concert.TicketsLeft == 0 {
		concert.TicketsLeft = concert.TicketsTotal
	}

	err = h.Service.CreateConcert(&concert)
	if err != nil {
		http.Error(w, "cannot create concert", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(concert)
}

// BuyTicket godoc
// @Summary Купить билет
// @Description Создаёт запрос на билет и уменьшает tickets_left
// @Tags concerts
// @Produce json
// @Security BearerAuth
// @Param id path int true "Concert ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid concert id"
// @Failure 401 {string} string "missing authorization header"
// @Failure 400 {string} string "no tickets available or ticket cannot be created"
// @Router /concerts/{id}/buy-ticket [post]
func (h *ConcertHandler) BuyTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/concerts/")
	path = strings.TrimSuffix(path, "/buy-ticket")

	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "invalid concert id", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value(UserIDKey)
	if userIDValue == nil {
		http.Error(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := userIDValue.(float64)
	if !ok {
		http.Error(w, "invalid user_id type", http.StatusUnauthorized)
		return
	}

	userID := int(userIDFloat)

	err = h.Service.BuyTicket(id, userID)
	if err != nil {
		http.Error(w, "no tickets available or ticket cannot be created", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "ticket request created successfully",
	})
}
