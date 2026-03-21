package model

type TicketEvent struct {
	UserID    int `json:"user_id"`
	ConcertID int `json:"concert_id"`
}
