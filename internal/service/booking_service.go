package service

import (
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

type BookingService struct{}

func NewBookingService() *BookingService {
	return &BookingService{}
}
func (s *BookingService) Create(b model.Booking) error {
	b.ID = repository.NextBookingID()
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return repository.CreateBooking(b)
}

func (s *BookingService) GetAll() []model.Booking {
	return repository.GetBookings()
}

func (s *BookingService) GetByID(id int) (model.Booking, bool) {
	return repository.GetBookingById(id)
}

func (s *BookingService) Update(id int, upd model.Booking) (model.Booking, bool, error) {
	// сохраним CreatedAt от старой записи, если есть
	old, ok := repository.GetBookingById(id)
	if !ok {
		return model.Booking{}, false, nil
	}
	upd.ID = id
	upd.CreatedAt = old.CreatedAt
	upd.UpdatedAt = time.Now()

	return repository.UpdateBookingById(id, upd)
}

func (s *BookingService) Delete(id int) (bool, error) {
	return repository.DeleteBookingById(id)
}
