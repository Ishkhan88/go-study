package service

import "github.com/Ishkhan88/go-study/internal/repository/postgres"

func InitPostgresRepos(concerts *postgres.ConcertsRepository, bookings *postgres.BookingsRepository) {
	InitPostgresConcertRepo(concerts)
	InitPostgresBookingRepo(bookings)
}
