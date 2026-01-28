package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/user" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		UserCreate(w, r)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/user/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		UserGet(w, r, id)
	case http.MethodPut:
		UserUpdate(w, r, id)
	case http.MethodDelete:
		UserDelete(w, r, id)
	default:
		writeErr(w, 405, "method not allowed")
	}
}

func ConcertHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/concert" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		ConcertCreate(w, r)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/concert/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		ConcertGet(w, r, id)
	case http.MethodPut:
		ConcertUpdate(w, r, id)
	case http.MethodDelete:
		ConcertDelete(w, r, id)
	default:
		writeErr(w, 405, "method not allowed")
	}
}

func BookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/booking" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		BookingCreate(w, r)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/booking/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		BookingGet(w, r, id)
	case http.MethodPut:
		BookingUpdate(w, r, id)
	case http.MethodDelete:
		BookingDelete(w, r, id)
	default:
		writeErr(w, 405, "method not allowed")
	}
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/notification" {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "only POST")
			return
		}
		NotificationCreate(w, r)
		return
	}

	id, ok := parseID(r.URL.Path, "/api/notification/")
	if !ok {
		writeErr(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		NotificationGet(w, r, id)
	case http.MethodPut:
		NotificationUpdate(w, r, id)
	case http.MethodDelete:
		NotificationDelete(w, r, id)
	default:
		writeErr(w, 405, "method not allowed")
	}
}
