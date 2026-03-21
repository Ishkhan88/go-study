package model

import "time"

type Notification struct {
	ID        int
	ConcertID int
	UserID    int
	Status    string
	SentAt    time.Time
}
