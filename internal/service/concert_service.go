package service

import (
	"context"
	"errors"

	"github.com/Ishkhan88/go-study/internal/core/port/clock"
	"github.com/Ishkhan88/go-study/internal/core/usecase"
	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
	"github.com/Ishkhan88/go-study/internal/repository/postgres"
)

var pgConcertsRepo *postgres.ConcertsRepository

func InitPostgresConcertRepo(repo *postgres.ConcertsRepository) {
	pgConcertsRepo = repo
}

func concertUC() usecase.ConcertUsecase {
	return usecase.NewConcertUsecase(
		repository.ConcertRepoAdapter{},
		clock.SystemClock{},
	)
}

func CreateConcert(c model.Concert) (model.Concert, error) {
	if pgConcertsRepo != nil {
		ctx := context.Background()

		created, err := pgConcertsRepo.CreateConcert(
			ctx,
			c.Title,
			c.Date,
			int32(c.TicketsTotal),
			int32(c.TicketsLeft),
			nil,
			nil,
		)
		if err != nil {
			return model.Concert{}, err
		}

		out := model.Concert{
			ID:             int(created.ID),
			Title:          created.Title,
			Date:           created.Date,
			Location:       c.Location,
			TicketPrice:    c.TicketPrice,
			TicketsTotal:   int(created.TicketsTotal),
			TicketsLeft:    int(created.TicketsLeft),
			OrganizerEmail: c.OrganizerEmail,
			CreatedAt:      created.CreatedAt,
			UpdatedAt:      created.CreatedAt,
		}
		return out, nil
	}

	return concertUC().Create(context.Background(), c)
}

func GetConcert(id int) (model.Concert, error) {
	if pgConcertsRepo != nil {
		got, err := pgConcertsRepo.GetConcertByID(context.Background(), int64(id))
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				return model.Concert{}, ErrNotFound
			}
			return model.Concert{}, err
		}

		return model.Concert{
			ID:           int(got.ID),
			Title:        got.Title,
			Date:         got.Date,
			TicketsTotal: int(got.TicketsTotal),
			TicketsLeft:  int(got.TicketsLeft),
			CreatedAt:    got.CreatedAt,
			UpdatedAt:    got.CreatedAt,
		}, nil
	}

	return concertUC().Get(context.Background(), id)
}

func UpdateConcert(id int, upd model.Concert) (model.Concert, error) {
	return concertUC().Update(context.Background(), id, upd)
}

func DeleteConcert(id int) error {
	return concertUC().Delete(context.Background(), id)
}

func ListConcerts() []model.Concert {
	list, err := concertUC().List(context.Background())
	if err != nil {
		return []model.Concert{}
	}
	return list
}
