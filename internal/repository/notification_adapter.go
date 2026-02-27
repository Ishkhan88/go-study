package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Ishkhan88/go-study/internal/core/port"
	"github.com/Ishkhan88/go-study/internal/model"
)

type NotificationRepoAdapter struct{}

var _ port.NotificationRepository = (*NotificationRepoAdapter)(nil)

func (NotificationRepoAdapter) NextID(ctx context.Context) (int, error) {
	_ = ctx

	muNotification.Lock()
	defer muNotification.Unlock()

	maxID := 0
	for i := range notification {
		if notification[i].ID > maxID {
			maxID = notification[i].ID
		}
	}
	return maxID + 1, nil
}

func (NotificationRepoAdapter) Add(ctx context.Context, n model.Notification) error {
	_ = ctx
	return addNotification(n)
}

func (NotificationRepoAdapter) GetByID(ctx context.Context, id int) (model.Notification, bool, error) {
	_ = ctx

	items := GetNotificationsSafeCopy()
	for _, n := range items {
		if n.ID == id {
			return n, true, nil
		}
	}
	return model.Notification{}, false, nil
}

func (NotificationRepoAdapter) Update(ctx context.Context, id int, upd model.Notification) (model.Notification, bool, error) {
	_ = ctx

	muNotification.Lock()
	defer muNotification.Unlock()

	idx := -1
	for i := range notification {
		if notification[i].ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return model.Notification{}, false, nil
	}

	upd.ID = id
	notification[idx] = upd

	if err := overwriteNotificationsLocked(); err != nil {
		return model.Notification{}, false, err
	}

	return notification[idx], true, nil
}

func (NotificationRepoAdapter) Delete(ctx context.Context, id int) (bool, error) {
	_ = ctx

	muNotification.Lock()
	defer muNotification.Unlock()

	idx := -1
	for i := range notification {
		if notification[i].ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false, nil
	}

	notification = append(notification[:idx], notification[idx+1:]...)

	if err := overwriteNotificationsLocked(); err != nil {
		return false, err
	}

	return true, nil
}

func (NotificationRepoAdapter) List(ctx context.Context) ([]model.Notification, error) {
	_ = ctx
	return GetNotificationsSafeCopy(), nil
}

func overwriteNotificationsLocked() error {
	f, err := os.OpenFile(notificationsFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open notifications file: %w", err)
	}
	defer f.Close()

	for _, n := range notification {
		raw, err := json.Marshal(n)
		if err != nil {
			return fmt.Errorf("marshal notification: %w", err)
		}
		if _, err := f.Write(append(raw, '\n')); err != nil {
			return fmt.Errorf("write notification line: %w", err)
		}
	}
	return nil
}
