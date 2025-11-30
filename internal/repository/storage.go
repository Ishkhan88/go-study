package repository

import (
	"fmt"

	"github.com/Ishkhan88/go-study/internal/model"
)

var (
	booking      []model.Booking
	concert      []model.Concert
	user         []model.User
	notification []model.Notification
)

// SaveEntity принимает любую сущность, которая реализует интерфейс model.Entity,
// определяет её реальный тип и кладёт в нужный слайс.
func SaveEntity(e model.Entity) {
	switch v := e.(type) {
	case model.User:
		user = append(user, v)
		fmt.Printf("Добавили пользователя: ID=%d\n", v.GetID())

	case model.Concert:
		concert = append(concert, v)
		fmt.Printf("Добавили концерт: ID=%d\n", v.GetID())

	case model.Booking:
		booking = append(booking, v)
		fmt.Printf("Добавили бронирование: ID=%d\n", v.GetID())

	case model.Notification:
		notification = append(notification, v)
		fmt.Printf("Добавили уведомление: ID=%d\n", v.GetID())

	default:
		// на случай, если кто-то передаст тип, который не учли
		fmt.Printf("Неизвестный тип сущности: %T\n", v)
	}
}

// чтобы можно было посмотреть, что реально накопилось в "коробках".
func GetUser() []model.User {
	return user
}

func GetConcert() []model.Concert {
	return concert
}

func GetBooking() []model.Booking {
	return booking
}

func GetNotification() []model.Notification {
	return notification
}
