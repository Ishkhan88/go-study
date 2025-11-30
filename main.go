package main

import (
	"fmt"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
	"github.com/Ishkhan88/go-study/internal/service"
)

func main() {
	// Пользователь (публичная модель)
	u := model.User{
		ID:        1,
		FirstName: "Ishkhan",
		LastName:  "Khalatyan",
		Email:     "Khalatian88@yandex.ru",
		Phone:     "+995500112233",
		AvatarURL: "https://example.com/avatar.png",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Данные для аутентификации (приватные поля через методы)
	//creds := auth.NewCredentials(u.ID, "<bcrypt-hash>")
	//creds.SetLastLoginAt(time.Now())

	// Концерт с ценой билета
	c := model.Concert{
		ID:             10,
		Title:          "Symphonic Rock Night",
		Date:           time.Now().Add(7 * 24 * time.Hour),
		Location:       "Moscow Arena",
		TicketPrice:    1500.99,
		TicketsTotal:   100,
		TicketsLeft:    100,
		OrganizerEmail: "organizer@example.com",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Бронь (используем константы статусов)
	b := model.Booking{
		ID:        1000,
		UserID:    u.ID,
		ConcertID: c.ID,
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Запись об уведомлении
	n := model.Notification{
		ID:        5000,
		UserID:    u.ID,
		ConcertID: c.ID,
		Status:    model.StatusSuccess,
		SentAt:    time.Now(),
	}

	fmt.Printf("User: %s %s, email=%s, phone=%s\n", u.FirstName, u.LastName, u.Email, u.Phone)
	//fmt.Println("Auth (privates via methods): userID=", creds.UserID(), "hash=", creds.PasswordHash())
	fmt.Printf("Concert: %s @ %s, price=%.2f, left=%d\n", c.Title, c.Location, c.TicketPrice, c.TicketsLeft)
	fmt.Printf("Booking: id=%d, status=%s\n", b.ID, b.Status)
	fmt.Printf("Notification: status=%s at %s\n", n.Status, n.SentAt.Format("2006-01-02 15:04:05"))

	ch := make(chan model.Entity)

	go repository.StartSaver(ch)
	go repository.NewItemsLogger(200 * time.Millisecond)
	go service.StartGenerator(ch, 2*time.Second)

	select {}
}
