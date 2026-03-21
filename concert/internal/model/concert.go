package model

import "time"

type Concert struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Location       string    `json:"location"`
	Date           time.Time `json:"date"`
	TicketsTotal   int       `json:"tickets_total"`
	TicketsLeft    int       `json:"tickets_left"`
	OrganizerEmail string    `json:"organizer_email"`
}
