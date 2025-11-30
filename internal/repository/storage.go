package repository

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
)

var (
	booking      []model.Booking
	concert      []model.Concert
	user         []model.User
	notification []model.Notification

	muUser         sync.Mutex
	muConcert      sync.Mutex
	muBooking      sync.Mutex
	muNotification sync.Mutex
)

// SaveEntity принимает любую сущность, которая реализует интерфейс model.Entity,
// определяет её реальный тип и кладёт в нужный слайс.
func SaveEntity(e model.Entity) {
	switch v := e.(type) {
	case model.User:
		muUser.Lock()
		user = append(user, v)
		muUser.Unlock()
		fmt.Printf("Добавили пользователя: ID=%d\n", v.GetID())

	case model.Concert:
		muConcert.Lock()
		concert = append(concert, v)
		muConcert.Unlock()
		fmt.Printf("Добавили концерт: ID=%d\n", v.GetID())

	case model.Booking:
		muBooking.Lock()
		booking = append(booking, v)
		muBooking.Unlock()
		fmt.Printf("Добавили бронирование: ID=%d\n", v.GetID())

	case model.Notification:
		muNotification.Lock()
		notification = append(notification, v)
		muNotification.Unlock()
		fmt.Printf("Добавили уведомление: ID=%d\n", v.GetID())

	default:
		// на случай, если кто-то передаст тип, который не учли
		log.Printf("Неизвестный тип сущности: %T\n", v)
	}
}

// StartSaver - работает в горутине читает из канала и передаёт в SaveEntity
func StartSaver(ch <-chan model.Entity) {
	for e := range ch {
		SaveEntity(e)
	}
}

func GetUserSafeCopy() []model.User {
	muUser.Lock()
	defer muUser.Unlock()

	copySlice := make([]model.User, len(user))
	copy(copySlice, user)
	return copySlice
}

func GetConcertSafeCopy() []model.Concert {
	muConcert.Lock()
	defer muConcert.Unlock()

	copySlice := make([]model.Concert, len(concert))
	copy(copySlice, concert)
	return copySlice
}

func GetBookingSafeCopy() []model.Booking {
	muBooking.Lock()
	defer muBooking.Unlock()

	copySlice := make([]model.Booking, len(booking))
	copy(copySlice, booking)
	return copySlice
}

func GetNotificationSafeCopy() []model.Notification {
	muNotification.Lock()
	defer muNotification.Unlock()

	copySlice := make([]model.Notification, len(notification))
	copy(copySlice, notification)
	return copySlice
}

// NewItemsLogger — логгер, который проверяет новые элементы в слайсах и выводит их
func NewItemsLogger(interval time.Duration) {
	ticker := time.NewTicker(interval)

	var prevUser, prevConcert, prevBooking, prevNotification int

	for range ticker.C {

		u := GetUserSafeCopy()
		c := GetConcertSafeCopy()
		b := GetBookingSafeCopy()
		n := GetNotificationSafeCopy()

		if len(u) > prevUser {
			for i := prevUser; i < len(u); i++ {
				log.Printf("New User added: %+v\n", u[i])
			}
			prevUser = len(u)
		}

		if len(c) > prevConcert {
			for i := prevConcert; i < len(c); i++ {
				log.Printf("New Concert added: %+v\n", c[i])
			}
			prevConcert = len(c)
		}

		if len(b) > prevBooking {
			for i := prevBooking; i < len(b); i++ {
				log.Printf("New Booking added: %+v\n", b[i])
			}
			prevBooking = len(b)
		}

		if len(n) > prevNotification {
			for i := prevNotification; i < len(n); i++ {
				log.Printf("New Notification added: %+v\n", n[i])
			}
			prevNotification = len(n)
		}
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
