package service

import (
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func CreateNotification(n model.Notification) (model.Notification, error) {
	// Create не должен требовать ID от клиента
	if n.ConcertID == 0 || n.UserID == 0 {
		return model.Notification{}, ErrBadInput
	}

	// генерируем ID
	if n.ID == 0 {
		n.ID = repository.GetNextNotificationID()
	}

	if n.Status == "" {
		n.Status = model.StatusSuccess
	}

	if n.SentAt.IsZero() {
		n.SentAt = time.Now()
	}

	if err := repository.AddNotification(n); err != nil {
		return model.Notification{}, err
	}

	return n, nil
}

func GetNotification(id int) (model.Notification, error) {
	n, ok := repository.GetNotificationByID(id)
	if !ok {
		return model.Notification{}, ErrNotFound
	}
	return n, nil
}

func UpdateNotification(id int, upd model.Notification) (model.Notification, error) {
	old, ok := repository.GetNotificationByID(id)
	if !ok {
		return model.Notification{}, ErrNotFound
	}

	if upd.ConcertID == 0 || upd.UserID == 0 {
		return model.Notification{}, ErrBadInput
	}

	switch upd.Status {
	case model.StatusSuccess, model.StatusFailed:
	default:
		return model.Notification{}, ErrBadInput
	}

	upd.ID = id

	// если sent_at не передали — сохраняем старый
	if upd.SentAt.IsZero() {
		upd.SentAt = old.SentAt
	}

	updated, ok2, err := repository.UpdateNotification(id, upd)
	if err != nil {
		return model.Notification{}, err
	}
	if !ok2 {
		return model.Notification{}, ErrNotFound
	}

	return updated, nil
}

func DeleteNotification(id int) error {
	deleted, err := repository.DeleteNotification(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func ListNotifications() []model.Notification {
	return repository.GetNotificationSafeCopy()
}
