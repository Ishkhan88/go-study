package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Ishkhan88/go-study/internal/core/port"
	"github.com/Ishkhan88/go-study/internal/model"
)

// BookingRepoAdapter реализует port.BookingRepository поверх текущего хранилища repository
// (слайс booking + mutex + файл bookings.jsonl).
type BookingRepoAdapter struct{}

var _ port.BookingRepository = (*BookingRepoAdapter)(nil)

func (BookingRepoAdapter) NextID(ctx context.Context) (int, error) {
	_ = ctx

	muBooking.Lock()
	defer muBooking.Unlock()

	maxID := 0
	for i := range booking {
		if booking[i].ID > maxID {
			maxID = booking[i].ID
		}
	}
	return maxID + 1, nil
}

func (BookingRepoAdapter) Add(ctx context.Context, b model.Booking) error {
	_ = ctx
	// Мы внутри пакета repository, поэтому можем вызывать внутреннюю функцию addBooking.
	return addBooking(b)
}

func (BookingRepoAdapter) GetByID(ctx context.Context, id int) (model.Booking, bool, error) {
	_ = ctx

	// Берём безопасную копию и ищем в ней.
	items := GetBookingsSafeCopy()
	for _, b := range items {
		if b.ID == id {
			return b, true, nil
		}
	}
	return model.Booking{}, false, nil
}

func (BookingRepoAdapter) List(ctx context.Context) ([]model.Booking, error) {
	_ = ctx
	return GetBookingsSafeCopy(), nil
}

func (BookingRepoAdapter) Update(ctx context.Context, id int, upd model.Booking) (model.Booking, bool, error) {
	_ = ctx

	muBooking.Lock()
	defer muBooking.Unlock()

	idx := -1
	for i := range booking {
		if booking[i].ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return model.Booking{}, false, nil
	}

	// ID фиксируем по ключу обновления
	upd.ID = id
	booking[idx] = upd

	// Перезаписываем файл целиком актуальным состоянием
	if err := overwriteBookingsLocked(); err != nil {
		return model.Booking{}, false, err
	}

	return booking[idx], true, nil
}

func (BookingRepoAdapter) Delete(ctx context.Context, id int) (bool, error) {
	_ = ctx

	muBooking.Lock()
	defer muBooking.Unlock()

	idx := -1
	for i := range booking {
		if booking[i].ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false, nil
	}

	// удалить элемент из слайса
	booking = append(booking[:idx], booking[idx+1:]...)

	// Перезаписать файл целиком
	if err := overwriteBookingsLocked(); err != nil {
		return false, err
	}

	return true, nil
}

// overwriteBookingsLocked перезаписывает bookings.jsonl актуальным слайсом booking.
// IMPORTANT: вызывать ТОЛЬКО когда muBooking уже удерживается (Locked).
func overwriteBookingsLocked() error {
	f, err := os.OpenFile(bookingsFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open bookings file: %w", err)
	}
	defer f.Close()

	for _, b := range booking {
		raw, err := json.Marshal(b)
		if err != nil {
			return fmt.Errorf("marshal booking: %w", err)
		}
		if _, err := f.Write(append(raw, '\n')); err != nil {
			return fmt.Errorf("write booking line: %w", err)
		}
	}
	return nil
}
