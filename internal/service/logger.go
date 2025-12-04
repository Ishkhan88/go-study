package service

import (
	"context"
	"log"
	"time"

	"github.com/Ishkhan88/go-study/internal/repository"
)

// Глобальные переменные для хранения последних размеров (или лучше использовать структуру)
var (
	lastUsersCount         = 0
	lastConcertsCount      = 0
	lastBookingsCount      = 0
	lastNotificationsCount = 0
)

// NewItemsLogger запускает фоновую горутину, которая логирует новые элементы
func NewItemsLogger(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Логгер новых элементов остановлен")
				return
			case <-ticker.C:
				users := repository.GetUsersSafeCopy()
				concerts := repository.GetConcertsSafeCopy()
				bookings := repository.GetBookingsSafeCopy()
				notifications := repository.GetNotificationsSafeCopy()

				if len(users) > lastUsersCount {
					log.Println("Новые пользователи:", len(users)-lastUsersCount, "шт.")
					// Или: users[lastUsersCount:]
					lastUsersCount = len(users)
				}

				if len(concerts) > lastConcertsCount {
					log.Println("Новые концерты:", len(concerts)-lastConcertsCount, "шт.")
					lastConcertsCount = len(concerts)
				}

				if len(bookings) > lastBookingsCount {
					log.Println("Новые бронирования:", len(bookings)-lastBookingsCount, "шт.")
					lastBookingsCount = len(bookings)
				}

				if len(notifications) > lastNotificationsCount {
					log.Println("Новые уведомления:", len(notifications)-lastNotificationsCount, "шт.")
					lastNotificationsCount = len(notifications)
				}
			}
		}
	}()
}
