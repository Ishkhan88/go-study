package main

import (
	"log"

	"github.com/Ishkhan88/go-study/notification/internal/config"
	"github.com/Ishkhan88/go-study/notification/internal/consumer"
	"github.com/Ishkhan88/go-study/notification/internal/email"
	"github.com/Ishkhan88/go-study/notification/internal/repository"
)

func main() {
	cfg := config.Load()

	emailSender := email.NewSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)

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

	log.Println("Notification DB connected")

	notificationRepo := repository.NewNotificationRepository(db)

	consumerInstance, err := consumer.NewConsumer(
		"amqp://guest:guest@rabbitmq:5672/",
		notificationRepo,
		emailSender,
	)
	if err != nil {
		log.Fatal("rabbitmq connection error:", err)
	}

	err = consumerInstance.Consume("ticket_requests")
	if err != nil {
		log.Fatal(err)
	}
}
