package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
)

func StartGenerator(ctx context.Context, ch chan<- model.Entity, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var id int = 1

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Генератор остановлен по команде контекста.")
			return

		case <-ticker.C:

			user := model.User{
				ID:        id,
				FirstName: "John",
				LastName:  fmt.Sprintf("Mayer_%d", id),
				Email:     fmt.Sprintf("user%d@example.com", id),
				Phone:     "8-900-000-00-00",
			}
			ch <- user
			concert := model.Concert{
				Title:          "Rock Fest",
				Date:           time.Now().AddDate(0, 1, 0),
				Location:       "Moscow",
				TicketsTotal:   100,
				TicketsLeft:    100,
				TicketPrice:    1599.0,
				OrganizerEmail: "userOrganizer@example.com",
			}
			ch <- concert

			booking := model.Booking{
				ID:        id + 2000,
				UserID:    user.ID,
				ConcertID: concert.ID,
				Status:    model.StatusPending,
			}
			ch <- booking

			notification := model.Notification{
				ID:        id + 3000,
				UserID:    user.ID,
				ConcertID: concert.ID,
				Status:    "created",
			}
			ch <- notification
		}
	}
}
