package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func parseID(path, prefix string) (int, bool) {
	if !strings.HasPrefix(path, prefix) {
		return 0, false
	}
	idStr := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	id, err := strconv.Atoi(idStr)
	return id, err == nil
}

// ===================== USERS =====================

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	// GET /api/users
	if r.Method != http.MethodGet {
		writeErr(w, 405, "only GET")
		return
	}
	writeJSON(w, 200, repository.GetUserSafeCopy())
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	// POST /api/user
	// GET/PUT/DELETE /api/user/{id}
	if r.URL.Path == "/api/user" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		var u model.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeErr(w, 400, "bad json")
			return
		}
		if u.ID == 0 || u.FirstName == "" || u.Email == "" {
			writeErr(w, 400, "required: id, first_name, email")
			return
		}
		now := time.Now()
		if u.CreatedAt.IsZero() {
			u.CreatedAt = now
		}
		u.UpdatedAt = now

		if err := repository.AddUser(u); err != nil {
			writeErr(w, 400, err.Error())
			return
		}
		writeJSON(w, 201, u)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/user/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		u, found := repository.GetUserByID(id)
		if !found {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, u)

	case http.MethodPut:
		var upd model.User
		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
			writeErr(w, 400, "bad json")
			return
		}
		if upd.FirstName == "" || upd.Email == "" {
			writeErr(w, 400, "required: first_name, email")
			return
		}
		old, found := repository.GetUserByID(id)
		if !found {
			writeErr(w, 404, "not found")
			return
		}
		upd.CreatedAt = old.CreatedAt
		upd.UpdatedAt = time.Now()
		upd.ID = id

		_, ok2, err := repository.UpdateUser(id, upd)
		if err != nil {
			writeErr(w, 500, "update failed")
			return
		}
		if !ok2 {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, upd)

	case http.MethodDelete:
		deleted, err := repository.DeleteUser(id)
		if err != nil {
			writeErr(w, 500, "delete failed")
			return
		}
		if !deleted {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, map[string]any{"deleted": true, "id": id})

	default:
		writeErr(w, 405, "method not allowed")
	}
}

// ===================== CONCERTS =====================

func ConcertsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, 405, "only GET")
		return
	}
	writeJSON(w, 200, repository.GetConcertSafeCopy())
}

func ConcertHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/concert" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		var c model.Concert
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			writeErr(w, 400, "bad json")
			return
		}
		if c.ID == 0 || c.Title == "" || c.Location == "" || c.OrganizerEmail == "" || c.TicketsTotal <= 0 {
			writeErr(w, 400, "required: id, title, location, organizer_email, tickets_total>0")
			return
		}
		now := time.Now()
		if c.CreatedAt.IsZero() {
			c.CreatedAt = now
		}
		c.UpdatedAt = now
		if c.TicketsLeft == 0 {
			c.TicketsLeft = c.TicketsTotal
		}

		if err := repository.AddConcert(c); err != nil {
			writeErr(w, 400, err.Error())
			return
		}
		writeJSON(w, 201, c)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/concert/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		c, found := repository.GetConcertByID(id)
		if !found {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, c)

	case http.MethodPut:
		var upd model.Concert
		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
			writeErr(w, 400, "bad json")
			return
		}
		if upd.Title == "" || upd.Location == "" || upd.OrganizerEmail == "" || upd.TicketsTotal <= 0 {
			writeErr(w, 400, "required: title, location, organizer_email, tickets_total>0")
			return
		}
		old, found := repository.GetConcertByID(id)
		if !found {
			writeErr(w, 404, "not found")
			return
		}
		upd.ID = id
		upd.CreatedAt = old.CreatedAt
		upd.UpdatedAt = time.Now()

		_, ok2, err := repository.UpdateConcert(id, upd)
		if err != nil {
			writeErr(w, 500, "update failed")
			return
		}
		if !ok2 {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, upd)

	case http.MethodDelete:
		deleted, err := repository.DeleteConcert(id)
		if err != nil {
			writeErr(w, 500, "delete failed")
			return
		}
		if !deleted {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, map[string]any{"deleted": true, "id": id})

	default:
		writeErr(w, 405, "method not allowed")
	}
}

// ===================== BOOKINGS =====================

func BookingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, 405, "only GET")
		return
	}
	writeJSON(w, 200, repository.GetBookingSafeCopy())
}

func BookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/booking" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		var b model.Booking
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			writeErr(w, 400, "bad json")
			return
		}
		if b.ID == 0 || b.UserID == 0 || b.ConcertID == 0 || b.Status == "" {
			writeErr(w, 400, "required: id, user_id, concert_id, status")
			return
		}
		now := time.Now()
		if b.CreatedAt.IsZero() {
			b.CreatedAt = now
		}
		b.UpdatedAt = now

		if err := repository.AddBooking(b); err != nil {
			writeErr(w, 400, err.Error())
			return
		}
		writeJSON(w, 201, b)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/booking/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		b, found := repository.GetBookingByID(id)
		if !found {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, b)

	case http.MethodPut:
		var upd model.Booking
		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
			writeErr(w, 400, "bad json")
			return
		}
		if upd.UserID == 0 || upd.ConcertID == 0 || upd.Status == "" {
			writeErr(w, 400, "required: user_id, concert_id, status")
			return
		}
		old, found := repository.GetBookingByID(id)
		if !found {
			writeErr(w, 404, "not found")
			return
		}
		upd.ID = id
		upd.CreatedAt = old.CreatedAt
		upd.UpdatedAt = time.Now()

		_, ok2, err := repository.UpdateBooking(id, upd)
		if err != nil {
			writeErr(w, 500, "update failed")
			return
		}
		if !ok2 {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, upd)

	case http.MethodDelete:
		deleted, err := repository.DeleteBooking(id)
		if err != nil {
			writeErr(w, 500, "delete failed")
			return
		}
		if !deleted {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, map[string]any{"deleted": true, "id": id})

	default:
		writeErr(w, 405, "method not allowed")
	}
}

// ===================== NOTIFICATIONS =====================

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, 405, "only GET")
		return
	}
	writeJSON(w, 200, repository.GetNotificationSafeCopy())
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/notification" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		var n model.Notification
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			writeErr(w, 400, "bad json")
			return
		}
		if n.ID == 0 || n.UserID == 0 || n.ConcertID == 0 || n.Status == "" {
			writeErr(w, 400, "required: id, user_id, concert_id, status")
			return
		}
		if n.SentAt.IsZero() {
			n.SentAt = time.Now()
		}

		if err := repository.AddNotification(n); err != nil {
			writeErr(w, 400, err.Error())
			return
		}
		writeJSON(w, 201, n)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/notification/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		n, found := repository.GetNotificationByID(id)
		if !found {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, n)

	case http.MethodPut:
		var upd model.Notification
		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
			writeErr(w, 400, "bad json")
			return
		}
		if upd.UserID == 0 || upd.ConcertID == 0 || upd.Status == "" {
			writeErr(w, 400, "required: user_id, concert_id, status")
			return
		}
		upd.ID = id
		if upd.SentAt.IsZero() {
			upd.SentAt = time.Now()
		}

		_, ok2, err := repository.UpdateNotification(id, upd)
		if err != nil {
			writeErr(w, 500, "update failed")
			return
		}
		if !ok2 {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, upd)

	case http.MethodDelete:
		deleted, err := repository.DeleteNotification(id)
		if err != nil {
			writeErr(w, 500, "delete failed")
			return
		}
		if !deleted {
			writeErr(w, 404, "not found")
			return
		}
		writeJSON(w, 200, map[string]any{"deleted": true, "id": id})

	default:
		writeErr(w, 405, "method not allowed")
	}
}
