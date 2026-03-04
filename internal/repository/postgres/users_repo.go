package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type UsersRepository struct {
	db *DB
}

func NewUsersRepository(db *DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) CreateUser(ctx context.Context, name, email string) (User, error) {
	const q = `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at, is_active, role;
	`

	var u User
	err := r.db.Pool.QueryRow(ctx, q, name, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
		&u.IsActive,
		&u.Role,
	)
	return u, err
}

func (r *UsersRepository) GetUserByID(ctx context.Context, id int64) (User, error) {
	const q = `
		SELECT id, name, email, created_at, is_active, role
		FROM users
		WHERE id = $1;
	`

	var u User
	err := r.db.Pool.QueryRow(ctx, q, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
		&u.IsActive,
		&u.Role,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return u, err
}

func (r *UsersRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	const q = `
		SELECT id, name, email, created_at, is_active, role
		FROM users
		WHERE email = $1;
	`

	var u User
	err := r.db.Pool.QueryRow(ctx, q, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
		&u.IsActive,
		&u.Role,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return u, err
}

func (r *UsersRepository) ListUsers(ctx context.Context, limit, offset int32) ([]User, error) {
	const q = `
		SELECT id, name, email, created_at, is_active, role
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Pool.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]User, 0, limit)
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.CreatedAt,
			&u.IsActive,
			&u.Role,
		); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *UsersRepository) SetUserActive(ctx context.Context, id int64, isActive bool) error {
	const q = `
		UPDATE users
		SET is_active = $2
		WHERE id = $1;
	`

	ct, err := r.db.Pool.Exec(ctx, q, id, isActive)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
