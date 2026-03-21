package model

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	AvatarURL    string
}
