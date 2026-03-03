package repository

import "github.com/Ishkhan88/go-study/internal/model"

// CreateNotification добавляет уведомление
func CreateNotification(n model.Notification) error {
	return AddNotification(n)
}

// GetNotifications возвращает все уведомления
func GetNotifications() []model.Notification {
	return GetNotificationSafeCopy()
}

// GetNotificationById возвращает уведомление по id
func GetNotificationById(id int) (model.Notification, bool) {
	return GetNotificationByID(id)
}

// UpdateNotificationById обновляет уведомление по id
func UpdateNotificationById(id int, upd model.Notification) (model.Notification, bool, error) {
	return UpdateNotification(id, upd)
}

// DeleteNotificationById удаляет уведомление по id
func DeleteNotificationById(id int) (bool, error) {
	return DeleteNotification(id)
}

// NextNotificationID возвращает следующий ID уведомления
func NextNotificationID() int {
	return GetNextNotificationID()
}
