package postgres

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
	Role      string    `json:"role"`
}

type Concert struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Date         time.Time `json:"date"`
	TicketsTotal int32     `json:"tickets_total"`
	TicketsLeft  int32     `json:"tickets_left"`
	Genre        *string   `json:"genre,omitempty"`
	Description  *string   `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type Booking struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ConcertID int64     `json:"concert_id"`
	Status    string    `json:"status"`
	Quantity  int32     `json:"quantity"`
	Note      *string   `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
