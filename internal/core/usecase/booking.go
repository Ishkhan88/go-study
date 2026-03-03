package usecase

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/apperr"
	"github.com/Ishkhan88/go-study/internal/core/port"
	"github.com/Ishkhan88/go-study/internal/model"
)

type BookingUsecase struct {
	repo  port.BookingRepository
	clock port.Clock
}

func NewBookingUsecase(repo port.BookingRepository, clock port.Clock) BookingUsecase {
	return BookingUsecase{
		repo:  repo,
		clock: clock,
	}
}

func (uc BookingUsecase) Create(ctx context.Context, b model.Booking) (model.Booking, error) {
	// Базовая валидация входа
	if b.UserID <= 0 || b.ConcertID <= 0 {
		return model.Booking{}, apperr.ErrBadInput
	}

	id, err := uc.repo.NextID(ctx)
	if err != nil {
		return model.Booking{}, err
	}

	now := uc.clock.Now()

	b.ID = id
	if b.Status == "" {
		b.Status = model.StatusPending
	}

	b.CreatedAt = now
	b.UpdatedAt = now

	if err := uc.repo.Add(ctx, b); err != nil {
		return model.Booking{}, err
	}

	return b, nil
}

func (uc BookingUsecase) Get(ctx context.Context, id int) (model.Booking, error) {
	if id <= 0 {
		return model.Booking{}, apperr.ErrBadInput
	}

	b, ok, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return model.Booking{}, err
	}
	if !ok {
		return model.Booking{}, apperr.ErrNotFound
	}

	return b, nil
}

func (uc BookingUsecase) Update(ctx context.Context, id int, upd model.Booking) (model.Booking, error) {
	if id <= 0 {
		return model.Booking{}, apperr.ErrBadInput
	}

	// Обновляем updatedAt по бизнес-правилу
	upd.UpdatedAt = uc.clock.Now()

	b, ok, err := uc.repo.Update(ctx, id, upd)
	if err != nil {
		return model.Booking{}, err
	}
	if !ok {
		return model.Booking{}, apperr.ErrNotFound
	}

	return b, nil
}

func (uc BookingUsecase) Delete(ctx context.Context, id int) error {
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

func (uc BookingUsecase) List(ctx context.Context) ([]model.Booking, error) {
	return uc.repo.List(ctx)
}
