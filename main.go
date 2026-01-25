package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Ishkhan88/go-study/internal/config"
	"github.com/Ishkhan88/go-study/internal/repository"
	"github.com/Ishkhan88/go-study/internal/service"
)

func main() {
	cfg := config.Default()

	if err := repository.LoadFromFiles(); err != nil {
		log.Println("load error:", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/users", service.UsersHandler) // GET list
	mux.HandleFunc("/api/user", service.UserHandler)   // POST create
	mux.HandleFunc("/api/user/", service.UserHandler)  // GET/PUT/DELETE by id

	mux.HandleFunc("/api/concerts", service.ConcertsHandler)
	mux.HandleFunc("/api/concert", service.ConcertHandler)
	mux.HandleFunc("/api/concert/", service.ConcertHandler)

	mux.HandleFunc("/api/bookings", service.BookingsHandler)
	mux.HandleFunc("/api/booking", service.BookingHandler)
	mux.HandleFunc("/api/booking/", service.BookingHandler)

	mux.HandleFunc("/api/notifications", service.NotificationsHandler)
	mux.HandleFunc("/api/notification", service.NotificationHandler)
	mux.HandleFunc("/api/notification/", service.NotificationHandler)

	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: mux,
	}

	var wg sync.WaitGroup

	// 4) запускаем логгер (и ждём его завершение)
	wg.Add(1)
	go func() {
		defer wg.Done()
		service.NewItemsLogger(ctx, cfg.LogInterval)
	}()

	// 5) запускаем сервер (и ждём завершение горутины)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("API server started: http://localhost" + cfg.ServerAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("server error:", err)
			stop()
		}
		log.Println("HTTP server goroutine stopped")
	}()

	// 6) ждём сигнал
	<-ctx.Done()
	log.Println("Shutdown signal received...")

	// 7) корректно останавливаем http сервер
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("server shutdown error:", err)
	}

	// 8) ждём завершения логгера и сервера
	wg.Wait()
	log.Println("Graceful shutdown completed")
}
