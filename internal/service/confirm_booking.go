package service

import (
	"context"
	"errors"

	"github.com/Ishkhan88/go-study/internal/repository/postgres"
)

func ConfirmBooking(id int64) error {
	if id <= 0 {
		return ErrBadInput
	}
	if pgBookingsRepo == nil {
		return errors.New("postgres bookings repo is not initialized")
	}

	if err := pgBookingsRepo.ConfirmBookingTx(context.Background(), id); err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}
