package usecase

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/apperr"
	"github.com/Ishkhan88/go-study/internal/core/port"
	"github.com/Ishkhan88/go-study/internal/model"
)

type ConcertUsecase struct {
	repo  port.ConcertRepository
	clock port.Clock
}

func NewConcertUsecase(repo port.ConcertRepository, clock port.Clock) ConcertUsecase {
	return ConcertUsecase{repo: repo, clock: clock}
}

func (uc ConcertUsecase) Create(ctx context.Context, c model.Concert) (model.Concert, error) {
	// Create не должен требовать ID от клиента
	if c.Title == "" || c.Location == "" || c.OrganizerEmail == "" {
		return model.Concert{}, apperr.ErrBadInput
	}

	// если ID не передали — сгенерируем
	if c.ID == 0 {
		id, err := uc.repo.NextID(ctx)
		if err != nil {
			return model.Concert{}, err
		}
		c.ID = id
	}

	now := uc.clock.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	// если TicketsLeft не задан — логично приравнять к TicketsTotal
	if c.TicketsLeft == 0 && c.TicketsTotal > 0 {
		c.TicketsLeft = c.TicketsTotal
	}

	if err := uc.repo.Add(ctx, c); err != nil {
		return model.Concert{}, err
	}

	return c, nil
}

func (uc ConcertUsecase) Get(ctx context.Context, id int) (model.Concert, error) {
	if id <= 0 {
		return model.Concert{}, apperr.ErrBadInput
	}

	c, ok, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return model.Concert{}, err
	}
	if !ok {
		return model.Concert{}, apperr.ErrNotFound
	}
	return c, nil
}

func (uc ConcertUsecase) Update(ctx context.Context, id int, upd model.Concert) (model.Concert, error) {
	if id <= 0 {
		return model.Concert{}, apperr.ErrBadInput
	}

	old, ok, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return model.Concert{}, err
	}
	if !ok {
		return model.Concert{}, apperr.ErrNotFound
	}

	if upd.Title == "" || upd.Location == "" || upd.OrganizerEmail == "" {
		return model.Concert{}, apperr.ErrBadInput
	}

	upd.ID = id
	upd.CreatedAt = old.CreatedAt
	upd.UpdatedAt = uc.clock.Now()

	updated, ok2, err := uc.repo.Update(ctx, id, upd)
	if err != nil {
		return model.Concert{}, err
	}
	if !ok2 {
		return model.Concert{}, apperr.ErrNotFound
	}

	return updated, nil
}

func (uc ConcertUsecase) Delete(ctx context.Context, id int) error {
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

func (uc ConcertUsecase) List(ctx context.Context) ([]model.Concert, error) {
	return uc.repo.List(ctx)
}
