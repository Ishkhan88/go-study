package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Ishkhan88/go-study/notification/internal/email"
	"github.com/Ishkhan88/go-study/notification/internal/model"
	"github.com/Ishkhan88/go-study/notification/internal/repository"
	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	repo    *repository.NotificationRepository
	sender  *email.Sender
}

func NewConsumer(url string, repo *repository.NotificationRepository, sender *email.Sender) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		repo:    repo,
		sender:  sender,
	}, nil
}

func (c *Consumer) Consume(queue string) error {
	_, err := c.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Waiting for messages...")

	for msg := range msgs {
		var event model.TicketEvent

		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			log.Println("error decoding message:", err)
			continue
		}

		log.Printf("Received ticket request: user_id=%d concert_id=%d\n", event.UserID, event.ConcertID)

		to := "ishkhanmariya@gmail.com"

		subject := "New Ticket Request"

		body := fmt.Sprintf(
			"User %d requested a ticket for concert %d",
			event.UserID,
			event.ConcertID,
		)

		err = c.sender.Send(to, subject, body)
		if err != nil {
			log.Println("error sending email:", err)
		} else {
			log.Println("Email sent successfully")
		}

		notification := &model.Notification{
			ConcertID: event.ConcertID,
			UserID:    event.UserID,
			Status:    "sent",
		}

		err = c.repo.Create(notification)
		if err != nil {
			log.Println("error saving notification:", err)
			continue
		}

		log.Println("Notification saved to database")
	}

	return nil
}
