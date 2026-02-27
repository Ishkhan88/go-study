package port

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/model"
)

type UserRepository interface {
	NextID(ctx context.Context) (int, error)
	Add(ctx context.Context, u model.User) error
	GetByID(ctx context.Context, id int) (model.User, bool, error)
	Update(ctx context.Context, id int, upd model.User) (model.User, bool, error)
	Delete(ctx context.Context, id int) (bool, error)
	List(ctx context.Context) ([]model.User, error)
}
