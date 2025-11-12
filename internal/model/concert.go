package model

import "time"

// Concert — описание концерта.
type Concert struct {
	ID             int       `json:"id"`              // уникальный идентификатор концерта
	Title          string    `json:"title"`           // название концерта
	Date           time.Time `json:"date"`            // дата проведения
	Location       string    `json:"location"`        // место проведения
	TicketPrice    float64   `json:"ticket_price"`    // цена одного билета
	TicketsTotal   int       `json:"tickets_total"`   // общее количество билетов
	TicketsLeft    int       `json:"tickets_left"`    // оставшиеся билеты
	OrganizerEmail string    `json:"organizer_email"` // e-mail организатора
	CreatedAt      time.Time `json:"created_at"`      // дата создания записи
	UpdatedAt      time.Time `json:"updated_at"`      // дата обновления записи
}
