package notification

import "time"

// Notification — запись о попытке/факте отправки письма организатору.
type Notification struct {
 ID        int       json:"id"
 ConcertID int       json:"concert_id"
 UserID    int       json:"user_id"
 Status    string    json:"status"  // success/failed
 SentAt    time.Time json:"sent_at" // дата и время отправки
}