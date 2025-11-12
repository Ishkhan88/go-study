package model

import "time"

const (
	StatusSuccess = "success" //удачная отправка
	StatusFailed  = "failed"  //ошибка при отправке
)

// Notification — запись о попытке/факте отправки письма организатору.
type Notification struct {
	ID        int       `json:"id"`         // уникальный идентификатор уведомления
	ConcertID int       `json:"concert_id"` // ID концерта
	UserID    int       `json:"user_id"`    // ID пользователя
	Status    string    `json:"status"`     // success / failed
	SentAt    time.Time `json:"sent_at"`    // дата и время отправки
}
