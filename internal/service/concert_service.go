package service

import (
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func CreateConcert(c model.Concert) (model.Concert, error) {
	// Create не должен требовать ID от клиента
	if c.Title == "" || c.Location == "" || c.OrganizerEmail == "" {
		return model.Concert{}, ErrBadInput
	}

	// если ID не передали — сгенерируем
	if c.ID == 0 {
		c.ID = repository.GetNextConcertID()
	}

	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	// если TicketsLeft не задан — логично приравнять к TicketsTotal
	if c.TicketsLeft == 0 && c.TicketsTotal > 0 {
		c.TicketsLeft = c.TicketsTotal
	}

	if err := repository.AddConcert(c); err != nil {
		return model.Concert{}, err
	}

	return c, nil
}

func GetConcert(id int) (model.Concert, error) {
	c, ok := repository.GetConcertByID(id)
	if !ok {
		return model.Concert{}, ErrNotFound
	}
	return c, nil
}

func UpdateConcert(id int, upd model.Concert) (model.Concert, error) {
	old, ok := repository.GetConcertByID(id)
	if !ok {
		return model.Concert{}, ErrNotFound
	}

	if upd.Title == "" || upd.Location == "" || upd.OrganizerEmail == "" {
		return model.Concert{}, ErrBadInput
	}

	upd.ID = id
	upd.CreatedAt = old.CreatedAt
	upd.UpdatedAt = time.Now()

	updated, ok2, err := repository.UpdateConcert(id, upd)
	if err != nil {
		return model.Concert{}, err
	}
	if !ok2 {
		return model.Concert{}, ErrNotFound
	}

	return updated, nil
}

func DeleteConcert(id int) error {
	deleted, err := repository.DeleteConcert(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func ListConcerts() []model.Concert {
	return repository.GetConcertSafeCopy()
}
