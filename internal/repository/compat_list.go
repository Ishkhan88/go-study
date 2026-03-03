package repository

import "github.com/Ishkhan88/go-study/internal/model"

// Совместимость со старым кодом хендлеров (router_list.go)

func GetUsers() []model.User {
	return GetUserSafeCopy()
}

func GetConcerts() []model.Concert {
	return GetConcertSafeCopy()
}
