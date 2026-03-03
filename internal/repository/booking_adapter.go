package repository

import "github.com/Ishkhan88/go-study/internal/model"

// CreateBooking добавляет бронь
func CreateBooking(b model.Booking) error {
	return AddBooking(b)
}

// GetBookings возвращает все брони
func GetBookings() []model.Booking {
	return GetBookingSafeCopy()
}

// GetBookingById возвращает бронь по id
func GetBookingById(id int) (model.Booking, bool) {
	return GetBookingByID(id)
}

// UpdateBookingById обновляет бронь по id
func UpdateBookingById(id int, upd model.Booking) (model.Booking, bool, error) {
	return UpdateBooking(id, upd)
}

// DeleteBookingById удаляет бронь по id
func DeleteBookingById(id int) (bool, error) {
	return DeleteBooking(id)
}

// NextBookingID возвращает следующий ID брони
func NextBookingID() int {
	return GetNextBookingID()
}
