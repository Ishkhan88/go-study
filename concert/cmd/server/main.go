// @title Concert Service API
// @version 1.0
// @description API для концертов и покупки билетов
// @host localhost:8082
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/Ishkhan88/go-study/concert/docs"
	"github.com/Ishkhan88/go-study/concert/internal/config"
	"github.com/Ishkhan88/go-study/concert/internal/handler"
	"github.com/Ishkhan88/go-study/concert/internal/mq"
	"github.com/Ishkhan88/go-study/concert/internal/repository"
	"github.com/Ishkhan88/go-study/concert/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewPostgres(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	if err != nil {
		log.Fatal("database connection error:", err)
	}

	fmt.Println("Concert DB connected")

	publisher, err := mq.NewPublisher("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal("rabbitmq connection error:", err)
	}

	fmt.Println("RabbitMQ connected")

	concertRepo := repository.NewConcertRepository(db)
	concertService := service.NewConcertService(concertRepo, publisher)
	concertHandler := handler.NewConcertHandler(concertService)

	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Concert service is running")
	})

	mux.HandleFunc("/concerts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			concertHandler.GetConcerts(w, r)
		case http.MethodPost:
			concertHandler.CreateConcert(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/concerts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/buy-ticket") {
			handler.AuthMiddleware(concertHandler.BuyTicket)(w, r)
			return
		}

		if r.Method == http.MethodGet {
			concertHandler.GetConcertByID(w, r)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	fmt.Println("Concert service started on port 8082")

	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
