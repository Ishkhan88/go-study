package repository

import (
	"database/sql"

	"github.com/Ishkhan88/go-study/notification/internal/model"
)

type NotificationRepository struct {
	DB *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) Create(n *model.Notification) error {
	query := `
	INSERT INTO notifications (concert_id, user_id, status)
	VALUES ($1, $2, $3)
	`

	_, err := r.DB.Exec(query, n.ConcertID, n.UserID, n.Status)
	return err
}
