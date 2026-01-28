package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/Ishkhan88/go-study/docs"
	"github.com/Ishkhan88/go-study/internal/config"
	"github.com/Ishkhan88/go-study/internal/repository"
	"github.com/Ishkhan88/go-study/internal/service"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Go Study API
// @version         1.0
// @description     CRUD API for users, concerts, bookings, notifications.
// @host            localhost:8081
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization

func main() {
	_ = godotenv.Load()
	cfg := config.Default()

	if err := repository.LoadFromFiles(); err != nil {
		log.Println("load error:", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", service.LoginHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// USERS
	mux.HandleFunc("/api/users", service.UserHandler) // GET list
	mux.HandleFunc("/api/user", service.UserHandler)  // POST create
	mux.HandleFunc("/api/user/", service.UserHandler) // GET/PUT/DELETE by id

	// CONCERTS
	mux.HandleFunc("/api/concerts", service.ConcertHandler) // GET list
	mux.HandleFunc("/api/concert", service.ConcertHandler)  // POST create
	mux.HandleFunc("/api/concert/", service.ConcertHandler) // GET/PUT/DELETE by id

	// BOOKINGS
	mux.HandleFunc("/api/bookings", service.BookingHandler) // GET list
	mux.HandleFunc("/api/booking", service.BookingHandler)  // POST create
	mux.HandleFunc("/api/booking/", service.BookingHandler) // GET/PUT/DELETE by id

	// NOTIFICATIONS
	mux.HandleFunc("/api/notifications", service.NotificationHandler) // GET list
	mux.HandleFunc("/api/notification", service.NotificationHandler)  // POST create
	mux.HandleFunc("/api/notification/", service.NotificationHandler) // GET/PUT/DELETE by id

	server := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		service.NewItemsLogger(ctx, cfg.LogInterval)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("API server started on", cfg.ServerAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("server error:", err)
			stop()
		}
		log.Println("HTTP server goroutine stopped")
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("server shutdown error:", err)
	}

	wg.Wait()
	log.Println("Graceful shutdown completed")
}
