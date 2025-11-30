package repository

import (
	"context"
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
func StartSaver(ctx context.Context, wg *sync.WaitGroup, ch <-chan model.Entity) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return // exit goroutine
		case e := <-ch:
			SaveEntity(e)
		}
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
func NewItemsLogger(ctx context.Context, wg *sync.WaitGroup, interval time.Duration) {
	defer wg.Done()

	ticker := time.NewTicker(interval)

	var prevUser, prevConcert, prevBooking, prevNotification int

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			user := GetUserSafeCopy()
			concert := GetConcertSafeCopy()
			booking := GetBookingSafeCopy()
			notification := GetNotificationSafeCopy()

			if len(user) > prevUser {
				for i := prevUser; i < len(user); i++ {
					log.Printf("New User added: %+v\n", user[i])
				}
				prevUser = len(user)
			}

			if len(concert) > prevConcert {
				for i := prevConcert; i < len(concert); i++ {
					log.Printf("New Concert added: %+v\n", concert[i])
				}
				prevConcert = len(concert)
			}

			if len(booking) > prevBooking {
				for i := prevBooking; i < len(booking); i++ {
					log.Printf("New Booking added: %+v\n", booking[i])
				}
				prevBooking = len(booking)
			}

			if len(notification) > prevNotification {
				for i := prevNotification; i < len(notification); i++ {
					log.Printf("New Notification added: %+v\n", notification[i])
				}
				prevNotification = len(notification)
			}
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
