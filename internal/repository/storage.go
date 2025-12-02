package repository

import (
	"sync"

	"github.com/Ishkhan88/go-study/internal/model"
)

var (
	user         []model.User
	concert      []model.Concert
	booking      []model.Booking
	notification []model.Notification

	muUser         sync.Mutex
	muConcert      sync.Mutex
	muBooking      sync.Mutex
	muNotification sync.Mutex
)

func SaveEntity(e model.Entity) {
	switch v := e.(type) {

	case model.User:
		muUser.Lock()
		user = append(user, v)
		muUser.Unlock()

	case model.Concert:
		muConcert.Lock()
		concert = append(concert, v)
		muConcert.Unlock()

	case model.Booking:
		muBooking.Lock()
		booking = append(booking, v)
		muBooking.Unlock()

	case model.Notification:
		muNotification.Lock()
		notification = append(notification, v)
		muNotification.Unlock()
	}
}

// безопасные копии
func GetUsersSafeCopy() []model.User {
	muUser.Lock()
	defer muUser.Unlock()
	copySlice := make([]model.User, len(user))
	copy(copySlice, user)
	return copySlice
}

func GetConcertsSafeCopy() []model.Concert {
	muConcert.Lock()
	defer muConcert.Unlock()
	copySlice := make([]model.Concert, len(concert))
	copy(copySlice, concert)
	return copySlice
}

func GetBookingsSafeCopy() []model.Booking {
	muBooking.Lock()
	defer muBooking.Unlock()
	copySlice := make([]model.Booking, len(booking))
	copy(copySlice, booking)
	return copySlice
}

func GetNotificationsSafeCopy() []model.Notification {
	muNotification.Lock()
	defer muNotification.Unlock()
	copySlice := make([]model.Notification, len(notification))
	copy(copySlice, notification)
	return copySlice
}
