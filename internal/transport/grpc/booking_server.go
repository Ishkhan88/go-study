package grpc

import (
	"context"
	"errors"
	"log"
	"time"

	gostudyv1 "github.com/Ishkhan88/go-study/gen/gostudy/v1"
	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// --- helpers ---

func bookingToProto(b model.Booking) *gostudyv1.Booking {
	return &gostudyv1.Booking{
		Id:        int64(b.ID),
		UserId:    int64(b.UserID),
		ConcertId: int64(b.ConcertID),
		Status:    b.Status,
		CreatedAt: timestamppb.New(b.CreatedAt),
		UpdatedAt: timestamppb.New(b.UpdatedAt),
	}
}

func bookingFromProto(p *gostudyv1.Booking) model.Booking {
	if p == nil {
		return model.Booking{}
	}
	b := model.Booking{
		ID:        int(p.Id),
		UserID:    int(p.UserId),
		ConcertID: int(p.ConcertId),
		Status:    p.Status,
	}
	if p.CreatedAt != nil {
		b.CreatedAt = p.CreatedAt.AsTime()
	}
	if p.UpdatedAt != nil {
		b.UpdatedAt = p.UpdatedAt.AsTime()
	}
	return b
}

func mapServiceErrBooking(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, service.ErrBadInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

// --- BookingsService methods ---

func (s *BookingsServer) Create(ctx context.Context, req *gostudyv1.CreateBookingRequest) (*gostudyv1.BookingResponse, error) {
	b := bookingFromProto(req.GetBooking())

	created, err := service.CreateBooking(b)
	// После успешного создания брони — отправим уведомление в NotificationService.
	// Ошибка уведомления НЕ должна ломать бронь (просто логируем).
	addr := s.NotificationAddr
	if addr == "" {
		addr = "127.0.0.1:50052"
	}

	//go
	func(userID, concertID int) {
		ctxN, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(
			ctxN,
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			log.Printf("notify dial error: %v", err)
			return
		}
		defer conn.Close()

		cli := gostudyv1.NewNotificationsServiceClient(conn)
		_, err = cli.Create(ctxN, &gostudyv1.CreateNotificationRequest{
			Notification: &gostudyv1.Notification{
				UserId:    int64(userID),
				ConcertId: int64(concertID),
				Status:    "success",
			},
		})
		if err != nil {
			log.Printf("notify create error: %v", err)
		}
	}(created.UserID, created.ConcertID)
	if err != nil {
		return nil, mapServiceErrBooking(err)
	}

	return &gostudyv1.BookingResponse{Booking: bookingToProto(created)}, nil
}

func (s *BookingsServer) Get(ctx context.Context, req *gostudyv1.IdRequest) (*gostudyv1.BookingResponse, error) {
	got, err := service.GetBooking(int(req.GetId()))
	if err != nil {
		return nil, mapServiceErrBooking(err)
	}
	return &gostudyv1.BookingResponse{Booking: bookingToProto(got)}, nil
}

func (s *BookingsServer) Update(ctx context.Context, req *gostudyv1.UpdateBookingRequest) (*gostudyv1.BookingResponse, error) {
	upd := bookingFromProto(req.GetBooking())

	updated, err := service.UpdateBooking(int(req.GetId()), upd)
	if err != nil {
		return nil, mapServiceErrBooking(err)
	}

	return &gostudyv1.BookingResponse{Booking: bookingToProto(updated)}, nil
}

func (s *BookingsServer) Delete(ctx context.Context, req *gostudyv1.IdRequest) (*gostudyv1.DeleteResponse, error) {
	err := service.DeleteBooking(int(req.GetId()))
	if err != nil {
		return nil, mapServiceErrBooking(err)
	}
	return &gostudyv1.DeleteResponse{Deleted: true, Id: req.GetId()}, nil
}

func (s *BookingsServer) List(ctx context.Context, req *gostudyv1.ListRequest) (*gostudyv1.ListBookingsResponse, error) {
	bookings := service.ListBookings()
	out := make([]*gostudyv1.Booking, 0, len(bookings))
	for _, b := range bookings {
		out = append(out, bookingToProto(b))
	}
	return &gostudyv1.ListBookingsResponse{Bookings: out}, nil
}
