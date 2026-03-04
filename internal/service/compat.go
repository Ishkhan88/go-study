package service

import (
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

var bookingSvc = NewBookingService()
var notificationSvc = NewNotificationService()

func CreateBooking(b model.Booking) (model.Booking, error) {
	b.ID = repository.NextBookingID()
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now

	if err := repository.CreateBooking(b); err != nil {
		return model.Booking{}, err
	}
	return b, nil
}

func GetBooking(id int) (model.Booking, error) {
	b, ok := bookingSvc.GetByID(id)
	if !ok {
		return model.Booking{}, ErrNotFound
	}
	return b, nil
}

func UpdateBooking(id int, upd model.Booking) (model.Booking, error) {
	b, ok, err := bookingSvc.Update(id, upd)
	if err != nil {
		return model.Booking{}, err
	}
	if !ok {
		return model.Booking{}, ErrNotFound
	}
	return b, nil
}

func DeleteBooking(id int) error {
	ok, err := bookingSvc.Delete(id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func ListBookings() []model.Booking {
	return bookingSvc.GetAll()
}

// ---- Notification compat ----

func CreateNotification(n model.Notification) (model.Notification, error) {
	n.ID = repository.NextNotificationID()
	if n.SentAt.IsZero() {
		n.SentAt = time.Now()
	}

	if err := repository.CreateNotification(n); err != nil {
		return model.Notification{}, err
	}
	return n, nil
}

func GetNotification(id int) (model.Notification, error) {
	n, ok := notificationSvc.GetByID(id)
	if !ok {
		return model.Notification{}, ErrNotFound
	}
	return n, nil
}

func UpdateNotification(id int, upd model.Notification) (model.Notification, error) {
	n, ok, err := notificationSvc.Update(id, upd)
	if err != nil {
		return model.Notification{}, err
	}
	if !ok {
		return model.Notification{}, ErrNotFound
	}
	return n, nil
}

func DeleteNotification(id int) error {
	ok, err := notificationSvc.Delete(id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func ListNotifications() []model.Notification {
	return notificationSvc.GetAll()
}
