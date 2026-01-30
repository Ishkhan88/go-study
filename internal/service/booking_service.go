package service

import (
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func CreateBooking(b model.Booking) (model.Booking, error) {
	if b.ID == 0 || b.UserID == 0 || b.ConcertID == 0 {
		return model.Booking{}, ErrBadInput
	}

	if b.Status == "" {
		b.Status = model.StatusPending
	}

	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	b.UpdatedAt = now

	if err := repository.AddBooking(b); err != nil {
		return model.Booking{}, err
	}

	return b, nil
}

func GetBooking(id int) (model.Booking, error) {
	b, ok := repository.GetBookingByID(id)
	if !ok {
		return model.Booking{}, ErrNotFound
	}
	return b, nil
}

func UpdateBooking(id int, upd model.Booking) (model.Booking, error) {
	old, ok := repository.GetBookingByID(id)
	if !ok {
		return model.Booking{}, ErrNotFound
	}

	// обязательные поля
	if upd.UserID == 0 || upd.ConcertID == 0 {
		return model.Booking{}, ErrBadInput
	}

	// статус валидируем (по твоим константам)
	switch upd.Status {
	case model.StatusPending, model.StatusConfirmed, model.StatusRejected:
	default:
		return model.Booking{}, ErrBadInput
	}

	upd.ID = id
	upd.CreatedAt = old.CreatedAt
	upd.UpdatedAt = time.Now()

	updated, ok2, err := repository.UpdateBooking(id, upd)
	if err != nil {
		return model.Booking{}, err
	}
	if !ok2 {
		return model.Booking{}, ErrNotFound
	}

	return updated, nil
}

func DeleteBooking(id int) error {
	deleted, err := repository.DeleteBooking(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}
