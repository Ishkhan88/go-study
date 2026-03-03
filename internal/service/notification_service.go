package service

import (
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) Create(n model.Notification) error {
	n.ID = repository.NextNotificationID()
	if n.SentAt.IsZero() {
		n.SentAt = time.Now()
	}
	return repository.CreateNotification(n)
}

func (s *NotificationService) GetAll() []model.Notification {
	return repository.GetNotifications()
}

func (s *NotificationService) GetByID(id int) (model.Notification, bool) {
	return repository.GetNotificationById(id)
}

func (s *NotificationService) Update(id int, upd model.Notification) (model.Notification, bool, error) {
	upd.ID = id
	// если sentAt не передали — проставим текущее
	if upd.SentAt.IsZero() {
		upd.SentAt = time.Now()
	}
	return repository.UpdateNotificationById(id, upd)
}

func (s *NotificationService) Delete(id int) (bool, error) {
	return repository.DeleteNotificationById(id)
}
