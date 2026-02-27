package port

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/model"
)

type BookingRepository interface {
	NextID(ctx context.Context) (int, error)
	Add(ctx context.Context, b model.Booking) error
	GetByID(ctx context.Context, id int) (model.Booking, bool, error)
	Update(ctx context.Context, id int, upd model.Booking) (model.Booking, bool, error)
	Delete(ctx context.Context, id int) (bool, error)
	List(ctx context.Context) ([]model.Booking, error)
}
