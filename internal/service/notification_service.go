package service

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/core/port/clock"
	"github.com/Ishkhan88/go-study/internal/core/usecase"
	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func notificationUC() usecase.NotificationUsecase {
	return usecase.NewNotificationUsecase(
		repository.NotificationRepoAdapter{},
		clock.SystemClock{},
	)
}

func CreateNotification(n model.Notification) (model.Notification, error) {
	return notificationUC().Create(context.Background(), n)
}

func GetNotification(id int) (model.Notification, error) {
	return notificationUC().Get(context.Background(), id)
}

func UpdateNotification(id int, upd model.Notification) (model.Notification, error) {
	return notificationUC().Update(context.Background(), id, upd)
}

func DeleteNotification(id int) error {
	return notificationUC().Delete(context.Background(), id)
}

func ListNotifications() []model.Notification {
	list, err := notificationUC().List(context.Background())
	if err != nil {
		return []model.Notification{}
	}
	return list
}
