package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type ConcertsRepository struct {
	db *DB
}

func NewConcertsRepository(db *DB) *ConcertsRepository {
	return &ConcertsRepository{db: db}
}

func (r *ConcertsRepository) CreateConcert(
	ctx context.Context,
	title string,
	dateTime any, // можно time.Time, оставил any чтобы не упираться в твой слой моделей
	ticketsTotal, ticketsLeft int32,
	genre, description *string,
) (Concert, error) {
	const q = `
		INSERT INTO concerts (title, date, tickets_total, tickets_left, genre, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, date, tickets_total, tickets_left, genre, description, created_at;
	`

	var c Concert
	err := r.db.Pool.QueryRow(ctx, q, title, dateTime, ticketsTotal, ticketsLeft, genre, description).Scan(
		&c.ID,
		&c.Title,
		&c.Date,
		&c.TicketsTotal,
		&c.TicketsLeft,
		&c.Genre,
		&c.Description,
		&c.CreatedAt,
	)
	return c, err
}

func (r *ConcertsRepository) GetConcertByID(ctx context.Context, id int64) (Concert, error) {
	const q = `
		SELECT id, title, date, tickets_total, tickets_left, genre, description, created_at
		FROM concerts
		WHERE id = $1;
	`

	var c Concert
	err := r.db.Pool.QueryRow(ctx, q, id).Scan(
		&c.ID,
		&c.Title,
		&c.Date,
		&c.TicketsTotal,
		&c.TicketsLeft,
		&c.Genre,
		&c.Description,
		&c.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Concert{}, ErrNotFound
	}
	return c, err
}

func (r *ConcertsRepository) ListConcerts(ctx context.Context, limit, offset int32) ([]Concert, error) {
	const q = `
		SELECT id, title, date, tickets_total, tickets_left, genre, description, created_at
		FROM concerts
		ORDER BY date
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Pool.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Concert, 0, limit)
	for rows.Next() {
		var c Concert
		if err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Date,
			&c.TicketsTotal,
			&c.TicketsLeft,
			&c.Genre,
			&c.Description,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *ConcertsRepository) DecreaseTicketsLeft(ctx context.Context, concertID int64, qty int32) error {
	const q = `
		UPDATE concerts
		SET tickets_left = tickets_left - $2
		WHERE id = $1 AND tickets_left >= $2;
	`
	ct, err := r.db.Pool.Exec(ctx, q, concertID, qty)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("not enough tickets")
	}
	return nil
}
