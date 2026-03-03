package usecase

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/apperr"
	"github.com/Ishkhan88/go-study/internal/core/port"
	"github.com/Ishkhan88/go-study/internal/model"
)

type UserUsecase struct {
	repo  port.UserRepository
	clock port.Clock
}

func NewUserUsecase(repo port.UserRepository, clock port.Clock) UserUsecase {
	return UserUsecase{repo: repo, clock: clock}
}

func (uc UserUsecase) Create(ctx context.Context, u model.User) (model.User, error) {
	if u.FirstName == "" || u.Email == "" {
		return model.User{}, apperr.ErrBadInput
	}

	if u.ID == 0 {
		id, err := uc.repo.NextID(ctx)
		if err != nil {
			return model.User{}, err
		}
		u.ID = id
	}

	now := uc.clock.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now

	if err := uc.repo.Add(ctx, u); err != nil {
		return model.User{}, err
	}

	return u, nil
}

func (uc UserUsecase) Get(ctx context.Context, id int) (model.User, error) {
	if id <= 0 {
		return model.User{}, apperr.ErrBadInput
	}

	u, ok, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	if !ok {
		return model.User{}, apperr.ErrNotFound
	}
	return u, nil
}

func (uc UserUsecase) Update(ctx context.Context, id int, upd model.User) (model.User, error) {
	if id <= 0 {
		return model.User{}, apperr.ErrBadInput
	}

	old, ok, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	if !ok {
		return model.User{}, apperr.ErrNotFound
	}

	if upd.FirstName == "" || upd.Email == "" {
		return model.User{}, apperr.ErrBadInput
	}

	upd.ID = id
	upd.CreatedAt = old.CreatedAt
	upd.UpdatedAt = uc.clock.Now()

	updated, ok2, err := uc.repo.Update(ctx, id, upd)
	if err != nil {
		return model.User{}, err
	}
	if !ok2 {
		return model.User{}, apperr.ErrNotFound
	}

	return updated, nil
}

func (uc UserUsecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return apperr.ErrBadInput
	}

	ok, err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return apperr.ErrNotFound
	}
	return nil
}

func (uc UserUsecase) List(ctx context.Context) ([]model.User, error) {
	return uc.repo.List(ctx)
}
