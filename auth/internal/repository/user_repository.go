package repository

import (
	"database/sql"

	"github.com/Ishkhan88/go-study/auth/internal/model"
)

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `
	SELECT id, name, email, password_hash, COALESCE(avatar_url, '')
	FROM users
	WHERE email = $1
	`

	user := &model.User{}

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *model.User) error {

	query := `
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	return r.DB.QueryRow(
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID)
}
