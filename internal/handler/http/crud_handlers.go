package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/service"
)

// -------- USERS --------

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/user" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		if !requireAuth(w, r) {
			return
		}

		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeErr(w, 400, "bad json")
			return
		}

		created, err := service.CreateUser(u)
		if err != nil {
			mapServiceErr(w, err)
			return
		}

		writeJSON(w, 201, created)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/user/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		u, err := service.GetUser(id)
		if err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, u)

	case http.MethodPut:
		if !requireAuth(w, r) {
			return
		}

		var upd model.User
		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
			writeErr(w, 400, "bad json")
			return
		}

		updated, err := service.UpdateUser(id, upd)
		if err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, updated)

	case http.MethodDelete:
		if !requireAuth(w, r) {
			return
		}

		if err := service.DeleteUser(id); err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, map[string]any{"deleted": true, "id": id})

	default:
		writeErr(w, 405, "method not allowed")
	}
}

// -------- CONCERTS --------

func ConcertHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/concert" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		if !requireAuth(w, r) {
			return
		}

		var c model.Concert
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			writeErr(w, 400, "bad json")
			return
		}

		created, err := service.CreateConcert(c)
		if err != nil {
			mapServiceErr(w, err)
			return
		}

		writeJSON(w, 201, created)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/concert/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		c, err := service.GetConcert(id)
		if err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, c)

	case http.MethodPut:
		if !requireAuth(w, r) {
			return
		}

		var upd model.Concert
		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
			writeErr(w, 400, "bad json")
			return
		}

		updated, err := service.UpdateConcert(id, upd)
		if err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, updated)

	case http.MethodDelete:
		if !requireAuth(w, r) {
			return
		}

		if err := service.DeleteConcert(id); err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, map[string]any{"deleted": true, "id": id})

	default:
		writeErr(w, 405, "method not allowed")
	}
}

// -------- BOOKINGS --------

func BookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/booking" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		if !requireAuth(w, r) {
			return
		}

		var b model.Booking
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			writeErr(w, 400, "bad json")
			return
		}

		created, err := service.CreateBooking(b)
		if err != nil {
			mapServiceErr(w, err)
			return
		}

		writeJSON(w, 201, created)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/booking/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		b, err := service.GetBooking(id)
		if err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, b)

	case http.MethodPut:
		if !requireAuth(w, r) {
			return
		}

		var upd model.Booking
		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
			writeErr(w, 400, "bad json")
			return
		}

		updated, err := service.UpdateBooking(id, upd)
		if err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, updated)

	case http.MethodDelete:
		if !requireAuth(w, r) {
			return
		}

		if err := service.DeleteBooking(id); err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, map[string]any{"deleted": true, "id": id})

	default:
		writeErr(w, 405, "method not allowed")
	}
}

// -------- NOTIFICATIONS --------

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/notification" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		if !requireAuth(w, r) {
			return
		}

		var n model.Notification
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			writeErr(w, 400, "bad json")
			return
		}

		created, err := service.CreateNotification(n)
		if err != nil {
			mapServiceErr(w, err)
			return
		}

		writeJSON(w, 201, created)
		return
	}

	// /api/notification/{id}
	id, ok := parseID(r.URL.Path, "/api/notification/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		n, err := service.GetNotification(id)
		if err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, n)

	case http.MethodPut:
		if !requireAuth(w, r) {
			return
		}

		var upd model.Notification
		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
			writeErr(w, 400, "bad json")
			return
		}

		updated, err := service.UpdateNotification(id, upd)
		if err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, updated)

	case http.MethodDelete:
		if !requireAuth(w, r) {
			return
		}

		if err := service.DeleteNotification(id); err != nil {
			mapServiceErr(w, err)
			return
		}
		writeJSON(w, 200, map[string]any{"deleted": true, "id": id})

	default:
		writeErr(w, 405, "method not allowed")
	}
}

// -------- HELPERS --------

func mapServiceErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrBadInput):
		writeErr(w, 400, "bad request")
	case errors.Is(err, service.ErrNotFound):
		writeErr(w, 404, "not found")
	default:
		writeErr(w, 500, "internal error")
	}
}
