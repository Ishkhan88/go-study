package usecase

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/apperr"
	"github.com/Ishkhan88/go-study/internal/core/port"
	"github.com/Ishkhan88/go-study/internal/model"
)

type NotificationUsecase struct {
	repo  port.NotificationRepository
	clock port.Clock 
}

func NewNotificationUsecase(repo port.NotificationRepository, clock port.Clock) NotificationUsecase {
	return NotificationUsecase{repo: repo, clock: clock}
}

func (uc NotificationUsecase) Create(ctx context.Context, n model.Notification) (model.Notification, error) {
	// Минимальная валидация: если клиент передал ID, он должен быть положительным.
	// (Если ID = 0 — сгенерируем.)
	if n.ID < 0 {
		return model.Notification{}, apperr.ErrBadInput
	}

	if n.ID == 0 {
		id, err := uc.repo.NextID(ctx)
		if err != nil {
			return model.Notification{}, err
		}
		n.ID = id
	}

	if err := uc.repo.Add(ctx, n); err != nil {
		return model.Notification{}, err
	}

	return n, nil
}

func (uc NotificationUsecase) Get(ctx context.Context, id int) (model.Notification, error) {
	if id <= 0 {
		return model.Notification{}, apperr.ErrBadInput
	}

	n, ok, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return model.Notification{}, err
	}
	if !ok {
		return model.Notification{}, apperr.ErrNotFound
	}

	return n, nil
}

func (uc NotificationUsecase) Update(ctx context.Context, id int, upd model.Notification) (model.Notification, error) {
	if id <= 0 {
		return model.Notification{}, apperr.ErrBadInput
	}

		_, ok, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return model.Notification{}, err
	}
	if !ok {
		return model.Notification{}, apperr.ErrNotFound
	}

	// Фиксируем ID по ключу
	upd.ID = id

	n, ok2, err := uc.repo.Update(ctx, id, upd)
	if err != nil {
		return model.Notification{}, err
	}
	if !ok2 {
		return model.Notification{}, apperr.ErrNotFound
	}

	return n, nil
}

func (uc NotificationUsecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return apperr.ErrBadInput
	}

	ok, err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return apperr.ErrNotFound
	}

	return nil
}

func (uc NotificationUsecase) List(ctx context.Context) ([]model.Notification, error) {
	return uc.repo.List(ctx)
}
