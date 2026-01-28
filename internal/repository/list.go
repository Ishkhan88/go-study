package repository

import "github.com/Ishkhan88/go-study/internal/model"

// GetUsers возвращает всех пользователей.
func GetUsers() []model.User {
	out := make([]model.User, 0, len(users))
	for _, u := range users {
		out = append(out, u)
	}
	return out
}

// GetConcerts возвращает все концерты.
func GetConcerts() []model.Concert {
	out := make([]model.Concert, 0, len(concerts))
	for _, c := range concerts {
		out = append(out, c)
	}
	return out
}

// GetBookings возвращает все брони.
func GetBookings() []model.Booking {
	out := make([]model.Booking, 0, len(bookings))
	for _, b := range bookings {
		out = append(out, b)
	}
	return out
}

// GetNotifications возвращает все уведомления.
func GetNotifications() []model.Notification {
	out := make([]model.Notification, 0, len(notifications))
	for _, n := range notifications {
		out = append(out, n)
	}
	return out
}
