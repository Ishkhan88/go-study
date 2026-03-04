package service

import (
	"context"
	"errors"

	"github.com/Ishkhan88/go-study/internal/core/port/clock"
	"github.com/Ishkhan88/go-study/internal/core/usecase"
	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
	"github.com/Ishkhan88/go-study/internal/repository/postgres"
)

var pgBookingsRepo *postgres.BookingsRepository

func InitPostgresBookingRepo(repo *postgres.BookingsRepository) {
	pgBookingsRepo = repo
}

func bookingUC() usecase.BookingUsecase {
	return usecase.NewBookingUsecase(
		repository.BookingRepoAdapter{},
		clock.SystemClock{},
	)
}

func CreateBooking(b model.Booking) (model.Booking, error) {
	if pgBookingsRepo != nil {
		created, err := pgBookingsRepo.CreateBooking(
			context.Background(),
			int64(b.UserID),
			int64(b.ConcertID),
			1,
			nil,
		)
		if err != nil {
			return model.Booking{}, err
		}
		return model.Booking{
			ID:        int(created.ID),
			UserID:    int(created.UserID),
			ConcertID: int(created.ConcertID),
			Status:    created.Status,
			CreatedAt: created.CreatedAt,
			UpdatedAt: created.CreatedAt,
		}, nil
	}

	return bookingUC().Create(context.Background(), b)
}

func GetBooking(id int) (model.Booking, error) {
	if pgBookingsRepo != nil {
		got, err := pgBookingsRepo.GetBookingByID(context.Background(), int64(id))
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return model.Booking{}, ErrNotFound
			}
			return model.Booking{}, err
		}
		return model.Booking{
			ID:        int(got.ID),
			UserID:    int(got.UserID),
			ConcertID: int(got.ConcertID),
			Status:    got.Status,
			CreatedAt: got.CreatedAt,
			UpdatedAt: got.CreatedAt,
		}, nil
	}

	return bookingUC().Get(context.Background(), id)
}

func UpdateBooking(id int, upd model.Booking) (model.Booking, error) {

	return bookingUC().Update(context.Background(), id, upd)
}

func DeleteBooking(id int) error {

	return bookingUC().Delete(context.Background(), id)
}

func ListBookings() []model.Booking {
	list, err := bookingUC().List(context.Background())
	if err != nil {
		return []model.Booking{}
	}
	return list
}
