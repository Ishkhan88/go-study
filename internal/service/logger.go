package service

import (
	"context"
	"log"
	"time"

	"github.com/Ishkhan88/go-study/internal/repository"
)

func NewItemsLogger(ctx context.Context, interval time.Duration) {
	// 1) Baseline: сколько элементов уже есть на старте (после LoadFromFiles)
	lastUserCount := len(repository.GetUserSafeCopy())
	lastConcertCount := len(repository.GetConcertSafeCopy())
	lastBookingCount := len(repository.GetBookingSafeCopy())
	lastNotificationCount := len(repository.GetNotificationSafeCopy())

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("NewItemsLogger stopped")
			return

		case <-ticker.C:
			users := repository.GetUserSafeCopy()
			concerts := repository.GetConcertSafeCopy()
			bookings := repository.GetBookingSafeCopy()
			notifications := repository.GetNotificationSafeCopy()

			// 2) Если кто-то добавился — логируем только новые элементы
			if len(users) > lastUserCount {
				log.Printf("Added users: %v\n", users[lastUserCount:])
				lastUserCount = len(users)
			}

			if len(concerts) > lastConcertCount {
				log.Printf("Added concerts: %v\n", concerts[lastConcertCount:])
				lastConcertCount = len(concerts)
			}

			if len(bookings) > lastBookingCount {
				log.Printf("Added bookings: %v\n", bookings[lastBookingCount:])
				lastBookingCount = len(bookings)
			}

			if len(notifications) > lastNotificationCount {
				log.Printf("Added notifications: %v\n", notifications[lastNotificationCount:])
				lastNotificationCount = len(notifications)
			}
		}
	}
}
