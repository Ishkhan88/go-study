package model

import "time"

// Status значения можно оформить как константы.
const (
	StatusPending   = "pending"   //бронь в статусе "ожидание"
	StatusConfirmed = "confirmed" //бронь в статусе "подтверждено"
	StatusRejected  = "rejected"  //бронь в статусе "отклонено"
)

// Booking — бронь/покупка билета.
type Booking struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ConcertID int       `json:"concert_id"`
	Status    string    `json:"status"` // pending / confirmed / rejected
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b Booking) GetID() int {
	return b.ID
}
