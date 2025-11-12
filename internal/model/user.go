package model

import "time"

// User - содержит основные данные профиля.

type User struct {
	ID        int       `json:"id"`         // уникальный идентификатор пользователя
	FirstName string    `json:"first_name"` // имя
	LastName  string    `json:"last_name"`  // фамилия
	Email     string    `json:"email"`      // адрес электронной почты
	Phone     string    `json:"phone"`      // номер телефона
	AvatarURL string    `json:"avatar_url"` // ссылка на аватар
	CreatedAt time.Time `json:"created_at"` // когда создан профиль
	UpdatedAt time.Time `json:"updated_at"` // когда обновлён профиль
}
