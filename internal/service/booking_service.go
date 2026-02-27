package service

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/core/port/clock"
	"github.com/Ishkhan88/go-study/internal/core/usecase"
	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

// bookingUC собирает зависимости для booking usecase.
// Позже это лучше вынести в отдельный слой wiring (internal/app), но сейчас оставляем тут.
func bookingUC() usecase.BookingUsecase {
	return usecase.NewBookingUsecase(
		repository.BookingRepoAdapter{}, // адаптер репозитория под интерфейс ports.BookingRepository
		clock.SystemClock{},             // реальный clock
	)
}

// Ниже — публичный API service, который могли дергать HTTP/GRPC handlers.
// Он просто делегирует в usecase.

func CreateBooking(b model.Booking) (model.Booking, error) {
	return bookingUC().Create(context.Background(), b)
}

func GetBooking(id int) (model.Booking, error) {
	return bookingUC().Get(context.Background(), id)
}

func UpdateBooking(id int, upd model.Booking) (model.Booking, error) {
	return bookingUC().Update(context.Background(), id, upd)
}

func DeleteBooking(id int) error {
	return bookingUC().Delete(context.Background(), id)
}

// Если у тебя где-то ожидается именно []model.Booking без ошибки — оставляем так.
// Но правильнее было бы вернуть ([]model.Booking, error).
func ListBookings() []model.Booking {
	list, err := bookingUC().List(context.Background())
	if err != nil {
		return []model.Booking{}
	}
	return list
}
