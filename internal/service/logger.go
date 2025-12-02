package service

import (
	"log"
	"time"

	"github.com/Ishkhan88/go-study/internal/repository"
)

func NewItemsLogger(interval time.Duration) {
	var lastUsersCount, lastConcertsCount, lastBookingsCount, lastNotificationsCount int

	ticker := time.NewTicker(interval)

	for range ticker.C {
		users := repository.GetUsersSafeCopy()
		concerts := repository.GetConcertsSafeCopy()
		bookings := repository.GetBookingsSafeCopy()
		notifications := repository.GetNotificationsSafeCopy()

		if len(users) > lastUsersCount {
			log.Println("Новые пользователи:", users[lastUsersCount:])
			lastUsersCount = len(users)
		}

		if len(concerts) > lastConcertsCount {
			log.Println("Новые концерты:", concerts[lastConcertsCount:])
			lastConcertsCount = len(concerts)
		}

		if len(bookings) > lastBookingsCount {
			log.Println("Новые бронирования:", bookings[lastBookingsCount:])
			lastBookingsCount = len(bookings)
		}

		if len(notifications) > lastNotificationsCount {
			log.Println("Новые уведомления:", notifications[lastNotificationsCount:])
			lastNotificationsCount = len(notifications)
		}
	}
}
