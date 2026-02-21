package main

import (
	"context"
	"log"
	"time"

	gostudyv1 "github.com/Ishkhan88/go-study/gen/gostudy/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// Если запускаешь общий сервер (users+concerts) отдельно — укажи его адрес.
	// Можно запускать старый ./cmd/grpcserver на другом порту, например :50053.
	monolithAddr = "127.0.0.1:50053"

	bookingServiceAddr      = "127.0.0.1:50051"
	notificationServiceAddr = "127.0.0.1:50052"
)

func dial(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
}

func main() {
	// общий контекст на выполнение всего сценария
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ------------------------------------------------------------
	// Часть A. Подготовка данных (User + Concert)
	// Вариант 1: через монолитный grpcserver (если ты его запускаешь отдельно)
	// Вариант 2: задать userID/concertID вручную (если они уже есть в CSV)
	// ------------------------------------------------------------

	var userID int64
	var concertID int64

	// ===== ВАРИАНТ 1 (рекомендуется для чистого теста): монолитный сервер =====
	// Закомментируй этот блок, если не запускаешь монолит.
	{
		monoConn, err := dial(ctx, monolithAddr)
		if err != nil {
			log.Fatalf("monolith connect error (%s): %v", monolithAddr, err)
		}
		defer monoConn.Close()

		users := gostudyv1.NewUsersServiceClient(monoConn)
		concerts := gostudyv1.NewConcertsServiceClient(monoConn)

		// --- CREATE USER ---
		uCreateResp, err := users.Create(ctx, &gostudyv1.CreateUserRequest{
			User: &gostudyv1.User{
				FirstName: "Aram",
				LastName:  "Test",
				Email:     "a@b.com",
				Phone:     "+7920009009",
				AvatarUrl: "https://example.com/a.png",
			},
		})
		if err != nil {
			log.Fatalf("User Create error: %v", err)
		}
		u := uCreateResp.GetUser()
		userID = u.GetId()
		log.Printf("USER CREATED: id=%d name=%s %s", u.GetId(), u.GetFirstName(), u.GetLastName())

		// --- CREATE CONCERT ---
		cCreateResp, err := concerts.Create(ctx, &gostudyv1.CreateConcertRequest{
			Concert: &gostudyv1.Concert{
				Title:          "Rock Night",
				Location:       "Berlin",
				OrganizerEmail: "org@concert.com",
				TicketPrice:    19.99,
				TicketsTotal:   100,
			},
		})
		if err != nil {
			log.Fatalf("Concert Create error: %v", err)
		}
		c := cCreateResp.GetConcert()
		concertID = c.GetId()
		log.Printf("CONCERT CREATED: id=%d title=%s left=%d", c.GetId(), c.GetTitle(), c.GetTicketsLeft())
	}

	// ===== ВАРИАНТ 2: если монолит не запускаешь, просто задай ID вручную =====
	// userID = 1
	// concertID = 1

	if userID == 0 || concertID == 0 {
		log.Fatalf("userID/concertID are zero: userID=%d concertID=%d (set manually or run monolith)", userID, concertID)
	}

	// ------------------------------------------------------------
	// Часть B. Тест взаимодействия микросервисов
	// BookingService -> NotificationService
	// ------------------------------------------------------------

	// --- connect bookingserver ---
	bookingConn, err := dial(ctx, bookingServiceAddr)
	if err != nil {
		log.Fatalf("booking connect error (%s): %v", bookingServiceAddr, err)
	}
	defer bookingConn.Close()

	bookings := gostudyv1.NewBookingsServiceClient(bookingConn)

	// --- CREATE BOOKING (должно триггернуть авто-уведомление) ---
	bCreateResp, err := bookings.Create(ctx, &gostudyv1.CreateBookingRequest{
		Booking: &gostudyv1.Booking{
			UserId:    userID,
			ConcertId: concertID,
			// Status не задаём -> будет pending
		},
	})
	if err != nil {
		log.Fatalf("Booking Create error: %v", err)
	}
	b := bCreateResp.GetBooking()
	log.Printf("BOOKING CREATED: id=%d user=%d concert=%d status=%s", b.GetId(), b.GetUserId(), b.GetConcertId(), b.GetStatus())

	// Немного подождём, чтобы goroutine в bookingserver успела создать уведомление
	time.Sleep(200 * time.Millisecond)

	// --- connect notificationserver ---
	notifConn, err := dial(ctx, notificationServiceAddr)
	if err != nil {
		log.Fatalf("notification connect error (%s): %v", notificationServiceAddr, err)
	}
	defer notifConn.Close()

	notifications := gostudyv1.NewNotificationsServiceClient(notifConn)

	// --- LIST NOTIFICATIONS (должно быть >= 1 и содержать userID+concertID) ---
	nListResp, err := notifications.List(ctx, &gostudyv1.ListRequest{})
	if err != nil {
		log.Fatalf("Notification List error: %v", err)
	}
	log.Printf("NOTIFICATION LIST: total=%d", len(nListResp.GetNotifications()))

	// (по желанию) можно вывести последние 3 уведомления
	ns := nListResp.GetNotifications()
	start := 0
	if len(ns) > 3 {
		start = len(ns) - 3
	}
	for _, n := range ns[start:] {
		log.Printf("NOTIFICATION: id=%d user=%d concert=%d status=%s", n.GetId(), n.GetUserId(), n.GetConcertId(), n.GetStatus())
	}

	// --- CLEANUP (по желанию) ---
	// Можно удалить booking:
	// _, _ = bookings.Delete(ctx, &gostudyv1.IdRequest{Id: b.GetId()})
}
