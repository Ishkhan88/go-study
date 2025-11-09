package auth

import "time"

// Credentials — внутренние (секретные) данные для аутентификации.
// Поля хранятся приватно, наружу выдаём только через методы.
type Credentials struct {
	userID       int       // связь с users.User (по ID)
	passwordHash string    // приватно: хэш пароля (bcrypt, scrypt и т.п.)
	lastLoginAt  time.Time // приватно: последняя успешная аутентификация
}

// геттеры/сеттеры, показывающие идею инкапсуляции:
func (c *Credentials) UserID() int              { return c.userID }
func (c *Credentials) PasswordHash() string     { return c.passwordHash }
func (c *Credentials) SetPasswordHash(h string) { c.passwordHash = h }

func (c *Credentials) LastLoginAt() time.Time     { return c.lastLoginAt }
func (c *Credentials) SetLastLoginAt(t time.Time) { c.lastLoginAt = t }
