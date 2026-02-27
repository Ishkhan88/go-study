package repository

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/core/port"
	"github.com/Ishkhan88/go-study/internal/model"
)

type ConcertRepoAdapter struct{}

var _ port.ConcertRepository = (*ConcertRepoAdapter)(nil)

func (ConcertRepoAdapter) NextID(ctx context.Context) (int, error) {
	_ = ctx
	return GetNextConcertID(), nil
}

func (ConcertRepoAdapter) Add(ctx context.Context, c model.Concert) error {
	_ = ctx
	return AddConcert(c)
}

func (ConcertRepoAdapter) GetByID(ctx context.Context, id int) (model.Concert, bool, error) {
	_ = ctx
	c, ok := GetConcertByID(id)
	return c, ok, nil
}

func (ConcertRepoAdapter) Update(ctx context.Context, id int, upd model.Concert) (model.Concert, bool, error) {
	_ = ctx
	c, ok, err := UpdateConcert(id, upd)
	return c, ok, err
}

func (ConcertRepoAdapter) Delete(ctx context.Context, id int) (bool, error) {
	_ = ctx
	return DeleteConcert(id)
}

func (ConcertRepoAdapter) List(ctx context.Context) ([]model.Concert, error) {
	_ = ctx
	return GetConcertSafeCopy(), nil
}
