package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
)

// StartGenerator запускает периодическое создание данных и отправляет их в канал ch
func StartGenerator(ctx context.Context, ch chan<- model.Entity, interval time.Duration) {
	ticker := time.NewTicker(interval)
	// Обязательно останавливаем тикер при выходе из функции, чтобы избежать утечек
	defer ticker.Stop()

	// Переменная для генерации уникальных ID, инкрементируется при каждой итерации
	var id int = 1

	for {
		select {
		case <-ctx.Done():
			// Выход из горутины при отмене контекста
			fmt.Println("Генератор остановлен по команде контекста.")
			return

		case <-ticker.C:
			// Логика генерации данных (из вашего второго фрагмента)

			user := model.User{ // 1. Пользователь
				// Используем инкрементируемый ID для демонстрации уникальности
				ID:        id,
				FirstName: "John",
				LastName:  fmt.Sprintf("Mayer_%d", id), // Делаем фамилию уникальной
				Email:     fmt.Sprintf("user%d@example.com", id),
				Phone:     "8-900-000-00-00",
			}
			ch <- user // Отправляем пользователя в канал

			concert := model.Concert{ // 2. Концерт
				ID:             id + 1000,
				Title:          "Rock Fest",
				Date:           time.Now().AddDate(0, 1, 0), // через месяц
				Location:       "Moscow",
				TicketsTotal:   100,
				TicketsLeft:    100,
				TicketPrice:    1599.0,
				OrganizerEmail: "userOrganizer@example.com",
			}
			ch <- concert // Отправляем концерт в канал

			booking := model.Booking{ // 3. Бронирование
				ID:        id + 2000,
				UserID:    user.ID,
				ConcertID: concert.ID,
				// Предполагается, что model.StatusPending где-то определен (например, "pending")
				Status: model.StatusPending,
			}
			ch <- booking // Отправляем бронирование в канал

			notification := model.Notification{ // 4. Уведомление
				ID:        id + 3000,
				UserID:    user.ID,
				ConcertID: concert.ID,
				Status:    "created",
			}
			ch <- notification // Отправляем уведомление в канал

			// Инкрементируем ID для следующего цикла
			id++
		}
	}
}
