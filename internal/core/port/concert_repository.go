package port

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/model"
)

type ConcertRepository interface {
	NextID(ctx context.Context) (int, error)
	Add(ctx context.Context, c model.Concert) error
	GetByID(ctx context.Context, id int) (model.Concert, bool, error)
	Update(ctx context.Context, id int, upd model.Concert) (model.Concert, bool, error)
	Delete(ctx context.Context, id int) (bool, error)
	List(ctx context.Context) ([]model.Concert, error)
}
