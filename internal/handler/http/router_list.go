package http

import (
	"net/http"

	"github.com/Ishkhan88/go-study/internal/repository"
)

// UsersList godoc
// @Summary List users
// @Tags users
// @Produce json
// @Success 200 {array} model.User
// @Router /api/users [get]
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, 405, "only GET")
		return
	}
	writeJSON(w, 200, repository.GetUsers())
}

// ConcertsList godoc
// @Summary List concerts
// @Tags concerts
// @Produce json
// @Success 200 {array} model.Concert
// @Router /api/concerts [get]
func ConcertsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, 405, "only GET")
		return
	}
	writeJSON(w, 200, repository.GetConcerts())
}

// BookingsList godoc
// @Summary List bookings
// @Tags bookings
// @Produce json
// @Success 200 {array} model.Booking
// @Router /api/bookings [get]
func BookingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, 405, "only GET")
		return
	}
	writeJSON(w, 200, repository.GetBookings())
}

// NotificationsList godoc
// @Summary List notifications
// @Tags notifications
// @Produce json
// @Success 200 {array} model.Notification
// @Router /api/notifications [get]
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, 405, "only GET")
		return
	}
	writeJSON(w, 200, repository.GetNotifications())
}
