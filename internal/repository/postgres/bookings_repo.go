package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type BookingsRepository struct {
	db *DB
}

func NewBookingsRepository(db *DB) *BookingsRepository {
	return &BookingsRepository{db: db}
}

func (r *BookingsRepository) CreateBooking(ctx context.Context, userID, concertID int64, qty int32, note *string) (Booking, error) {
	const q = `
		INSERT INTO bookings (user_id, concert_id, status, quantity, note)
		VALUES ($1, $2, 'pending', $3, $4)
		RETURNING id, user_id, concert_id, status, quantity, note, created_at;
	`

	var b Booking
	err := r.db.Pool.QueryRow(ctx, q, userID, concertID, qty, note).Scan(
		&b.ID, &b.UserID, &b.ConcertID, &b.Status, &b.Quantity, &b.Note, &b.CreatedAt,
	)
	return b, err
}

func (r *BookingsRepository) GetBookingByID(ctx context.Context, id int64) (Booking, error) {
	const q = `
		SELECT id, user_id, concert_id, status, quantity, note, created_at
		FROM bookings
		WHERE id = $1;
	`

	var b Booking
	err := r.db.Pool.QueryRow(ctx, q, id).Scan(
		&b.ID, &b.UserID, &b.ConcertID, &b.Status, &b.Quantity, &b.Note, &b.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Booking{}, ErrNotFound
	}
	return b, err
}

func (r *BookingsRepository) ConfirmBookingTx(ctx context.Context, bookingID int64) error {
	return r.db.WithTx(ctx, func(tx pgx.Tx) error {
		const getBooking = `
			SELECT id, concert_id, quantity, status
			FROM bookings
			WHERE id = $1;
		`
		var (
			id        int64
			concertID int64
			qty       int32
			status    string
		)

		if err := tx.QueryRow(ctx, getBooking, bookingID).Scan(&id, &concertID, &qty, &status); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNotFound
			}
			return err
		}

		if status == "confirmed" {
			return nil
		}

		const decTickets = `
			UPDATE concerts
			SET tickets_left = tickets_left - $2
			WHERE id = $1 AND tickets_left >= $2;
		`
		ct, err := tx.Exec(ctx, decTickets, concertID, qty)
		if err != nil {
			return err
		}
		if ct.RowsAffected() == 0 {
			return errors.New("not enough tickets")
		}

		const setStatus = `
			UPDATE bookings
			SET status = 'confirmed'
			WHERE id = $1;
		`
		ct2, err := tx.Exec(ctx, setStatus, bookingID)
		if err != nil {
			return err
		}
		if ct2.RowsAffected() == 0 {
			return ErrNotFound
		}

		return nil
	})
}
