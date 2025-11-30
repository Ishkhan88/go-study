package service

import (
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
)

// StartGenerator запускает периодическое создание данных
func StartGenerator(ch chan<- model.Entity, interval time.Duration) {
	ticker := time.NewTicker(interval)
	id := 1

	for range ticker.C {

		user := model.User{ // 1. Пользователь
			ID:        123,
			FirstName: "John",
			LastName:  "Mayer",
			Email:     "user@example.com",
			Phone:     "8-900-000-00-00",
		}
		ch <- user

		concert := model.Concert{ // 2. Концерт
			ID:             345,
			Title:          "Rock Fest",
			Date:           time.Now().AddDate(0, 1, 0), // через месяц
			Location:       "Moscow",
			TicketsTotal:   100,
			TicketsLeft:    100,
			TicketPrice:    1599.0,
			OrganizerEmail: "userOrganizer@example.com",
		}
		ch <- concert

		booking := model.Booking{ // 3. Бронирование
			ID:        555,
			UserID:    user.ID,
			ConcertID: concert.ID,
			Status:    model.StatusPending,
		}
		ch <- booking

		notification := model.Notification{ // 4. Уведомление
			ID:        888,
			UserID:    user.ID,
			ConcertID: concert.ID,
			Status:    "created",
		}
		ch <- notification

		id++
	}
}
