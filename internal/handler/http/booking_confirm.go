package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ishkhan88/go-study/internal/repository/postgres"
)

type BookingConfirmHandler struct {
	Repo *postgres.BookingsRepository
}

func NewBookingConfirmHandler(repo *postgres.BookingsRepository) *BookingConfirmHandler {
	return &BookingConfirmHandler{Repo: repo}
}

func (h *BookingConfirmHandler) HandleConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 3 || parts[0] != "bookings" || parts[2] != "confirm" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid booking id"})
		return
	}

	err = h.Repo.ConfirmBookingTx(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == postgres.ErrNotFound {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "not enough tickets") {
			status = http.StatusConflict
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "confirmed"})
}
