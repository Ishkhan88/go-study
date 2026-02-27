package service

import (
	"context"

	"github.com/Ishkhan88/go-study/internal/core/port/clock"
	"github.com/Ishkhan88/go-study/internal/core/usecase"
	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func concertUC() usecase.ConcertUsecase {
	return usecase.NewConcertUsecase(
		repository.ConcertRepoAdapter{},
		clock.SystemClock{},
	)
}

func CreateConcert(c model.Concert) (model.Concert, error) {
	return concertUC().Create(context.Background(), c)
}

func GetConcert(id int) (model.Concert, error) {
	return concertUC().Get(context.Background(), id)
}

func UpdateConcert(id int, upd model.Concert) (model.Concert, error) {
	return concertUC().Update(context.Background(), id, upd)
}

func DeleteConcert(id int) error {
	return concertUC().Delete(context.Background(), id)
}

// Оставляем старую сигнатуру без error для совместимости с текущими хендлерами
func ListConcerts() []model.Concert {
	list, err := concertUC().List(context.Background())
	if err != nil {
		return []model.Concert{}
	}
	return list
}
