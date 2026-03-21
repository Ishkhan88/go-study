package repository

import (
	"database/sql"

	"github.com/Ishkhan88/go-study/concert/internal/model"
)

type ConcertRepository struct {
	DB *sql.DB
}

func NewConcertRepository(db *sql.DB) *ConcertRepository {
	return &ConcertRepository{DB: db}
}

func (r *ConcertRepository) GetAll() ([]model.Concert, error) {

	rows, err := r.DB.Query(`
	SELECT id, title, description, location, date, tickets_total, tickets_left, organizer_email
	FROM concerts
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var concerts []model.Concert

	for rows.Next() {
		var c model.Concert

		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Description,
			&c.Location,
			&c.Date,
			&c.TicketsTotal,
			&c.TicketsLeft,
			&c.OrganizerEmail,
		)
		if err != nil {
			return nil, err
		}

		concerts = append(concerts, c)
	}

	return concerts, nil
}

func (r *ConcertRepository) GetByID(id int) (*model.Concert, error) {

	query := `
	SELECT id, title, description, location, date, tickets_total, tickets_left, organizer_email
	FROM concerts
	WHERE id = $1
	`

	var c model.Concert

	err := r.DB.QueryRow(query, id).Scan(
		&c.ID,
		&c.Title,
		&c.Description,
		&c.Location,
		&c.Date,
		&c.TicketsTotal,
		&c.TicketsLeft,
		&c.OrganizerEmail,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *ConcertRepository) Create(concert *model.Concert) error {
	query := `
	INSERT INTO concerts (title, description, location, date, tickets_total, tickets_left, organizer_email)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	return r.DB.QueryRow(
		query,
		concert.Title,
		concert.Description,
		concert.Location,
		concert.Date,
		concert.TicketsTotal,
		concert.TicketsLeft,
		concert.OrganizerEmail,
	).Scan(&concert.ID)
}

func (r *ConcertRepository) DecreaseTickets(concertID int) error {
	query := `
	UPDATE concerts
	SET tickets_left = tickets_left - 1
	WHERE id = $1 AND tickets_left > 0
	`

	result, err := r.DB.Exec(query, concertID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ConcertRepository) CreateTicket(userID, concertID int) error {
	query := `
	INSERT INTO tickets (user_id, concert_id, status)
	VALUES ($1, $2, 'requested')
	`

	_, err := r.DB.Exec(query, userID, concertID)
	return err
}
