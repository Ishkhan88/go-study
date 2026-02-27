package repository

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/core/port"
	"github.com/Ishkhan88/go-study/internal/model"
)

type UserRepoAdapter struct{}

var _ port.UserRepository = (*UserRepoAdapter)(nil)

func (UserRepoAdapter) NextID(ctx context.Context) (int, error) {
	_ = ctx
	return GetNextUserID(), nil
}

func (UserRepoAdapter) Add(ctx context.Context, u model.User) error {
	_ = ctx
	return AddUser(u)
}

func (UserRepoAdapter) GetByID(ctx context.Context, id int) (model.User, bool, error) {
	_ = ctx
	u, ok := GetUserByID(id)
	return u, ok, nil
}

func (UserRepoAdapter) Update(ctx context.Context, id int, upd model.User) (model.User, bool, error) {
	_ = ctx
	return UpdateUser(id, upd)
}

func (UserRepoAdapter) Delete(ctx context.Context, id int) (bool, error) {
	_ = ctx
	return DeleteUser(id)
}

func (UserRepoAdapter) List(ctx context.Context) ([]model.User, error) {
	_ = ctx
	return GetUserSafeCopy(), nil
}
