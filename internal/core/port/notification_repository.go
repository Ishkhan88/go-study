package port

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/model"
)

type NotificationRepository interface {
	NextID(ctx context.Context) (int, error)
	Add(ctx context.Context, n model.Notification) error
	GetByID(ctx context.Context, id int) (model.Notification, bool, error)
	Update(ctx context.Context, id int, upd model.Notification) (model.Notification, bool, error)
	Delete(ctx context.Context, id int) (bool, error)
	List(ctx context.Context) ([]model.Notification, error)
}
