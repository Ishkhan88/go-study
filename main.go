package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ishkhan88/go-study/internal/repository"
	"github.com/Ishkhan88/go-study/internal/service"
)

func main() {
	// 1) Восстановление данных из CSV при старте (чтобы слайсы наполнились)
	if err := repository.LoadFromFiles(); err != nil {
		log.Println("load error:", err)
	}

	// 2) context для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 3) Роуты webserver
	mux := http.NewServeMux()

	// Users
	mux.HandleFunc("/api/users", service.UsersHandler) // GET list
	mux.HandleFunc("/api/user", service.UserHandler)   // POST create
	mux.HandleFunc("/api/user/", service.UserHandler)  // GET/PUT/DELETE by id

	// Concerts
	mux.HandleFunc("/api/concerts", service.ConcertsHandler)
	mux.HandleFunc("/api/concert", service.ConcertHandler)
	mux.HandleFunc("/api/concert/", service.ConcertHandler)

	// Bookings
	mux.HandleFunc("/api/bookings", service.BookingsHandler)
	mux.HandleFunc("/api/booking", service.BookingHandler)
	mux.HandleFunc("/api/booking/", service.BookingHandler)

	// Notifications
	mux.HandleFunc("/api/notifications", service.NotificationsHandler)
	mux.HandleFunc("/api/notification", service.NotificationHandler)
	mux.HandleFunc("/api/notification/", service.NotificationHandler)

	// 4) Запускаем логгер (после LoadFromFiles — чтобы НЕ логировать старые данные)
	go service.NewItemsLogger(ctx, 200*time.Millisecond)

	// 5) Запуск HTTP сервера
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("API server started: http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("server error:", err)
			stop() // если сервер упал — завершаем приложение
		}
	}()

	// 6) ждём сигнал ОС
	<-ctx.Done()
	log.Println("Shutdown signal received...")

	// 7) останавливаем сервер аккуратно (даем время завершить запросы)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("server shutdown error:", err)
	}

	log.Println("Graceful shutdown completed")
}
