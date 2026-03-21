package model

type Ticket struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ConcertID int    `json:"concert_id"`
	Status    string `json:"status"`
}
