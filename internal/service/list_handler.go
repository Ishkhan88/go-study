package service

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

	// Подстрой под реальные функции репозитория:
	// users := repository.GetUsers()
	// или users := repository.ListUsers()
	// или users := repository.UsersAll()

	users := repository.GetUsers() // <-- если такого нет, см. ниже
	writeJSON(w, 200, users)
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

	concerts := repository.GetConcerts() // <-- если такого нет, см. ниже
	writeJSON(w, 200, concerts)
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

	bookings := repository.GetBookings() // <-- если такого нет, см. ниже
	writeJSON(w, 200, bookings)
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

	notifications := repository.GetNotifications() // <-- если такого нет, см. ниже
	writeJSON(w, 200, notifications)
}
