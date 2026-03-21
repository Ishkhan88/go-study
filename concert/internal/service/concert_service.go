package service

import (
	"github.com/Ishkhan88/go-study/concert/internal/model"
	"github.com/Ishkhan88/go-study/concert/internal/mq"
	"github.com/Ishkhan88/go-study/concert/internal/repository"
)

type ConcertService struct {
	Repo      *repository.ConcertRepository
	Publisher *mq.Publisher
}

func NewConcertService(repo *repository.ConcertRepository, publisher *mq.Publisher) *ConcertService {
	return &ConcertService{
		Repo:      repo,
		Publisher: publisher,
	}
}

func (s *ConcertService) GetAllConcerts() ([]model.Concert, error) {
	return s.Repo.GetAll()
}

func (s *ConcertService) GetConcertByID(id int) (*model.Concert, error) {
	return s.Repo.GetByID(id)
}

func (s *ConcertService) CreateConcert(concert *model.Concert) error {
	return s.Repo.Create(concert)
}

func (s *ConcertService) BuyTicket(concertID, userID int) error {
	err := s.Repo.DecreaseTickets(concertID)
	if err != nil {
		return err
	}

	err = s.Repo.CreateTicket(userID, concertID)
	if err != nil {
		return err
	}

	event := model.TicketEvent{
		UserID:    userID,
		ConcertID: concertID,
	}

	err = s.Publisher.Publish("ticket_requests", event)
	if err != nil {
		return err
	}

	return nil
}
