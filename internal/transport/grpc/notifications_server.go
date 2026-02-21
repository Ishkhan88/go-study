package grpc

import (
	"context"
	"errors"

	gostudyv1 "github.com/Ishkhan88/go-study/gen/gostudy/v1"
	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// helpers

func notificationToProto(n model.Notification) *gostudyv1.Notification {
	return &gostudyv1.Notification{
		Id:        int64(n.ID),
		UserId:    int64(n.UserID),
		ConcertId: int64(n.ConcertID),
		Status:    n.Status,
		SentAt:    timestamppb.New(n.SentAt),
	}
}

func notificationFromProto(p *gostudyv1.Notification) model.Notification {
	if p == nil {
		return model.Notification{}
	}
	n := model.Notification{
		ID:        int(p.Id),
		UserID:    int(p.UserId),
		ConcertID: int(p.ConcertId),
		Status:    p.Status,
	}
	if p.SentAt != nil {
		n.SentAt = p.SentAt.AsTime()
	}
	return n
}

func mapServiceErrNotification(err error) error {
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

// methods

func (s *NotificationsServer) Create(ctx context.Context, req *gostudyv1.CreateNotificationRequest) (*gostudyv1.NotificationResponse, error) {
	n := notificationFromProto(req.GetNotification())

	created, err := service.CreateNotification(n)
	if err != nil {
		return nil, mapServiceErrNotification(err)
	}

	return &gostudyv1.NotificationResponse{Notification: notificationToProto(created)}, nil
}

func (s *NotificationsServer) Get(ctx context.Context, req *gostudyv1.IdRequest) (*gostudyv1.NotificationResponse, error) {
	got, err := service.GetNotification(int(req.GetId()))
	if err != nil {
		return nil, mapServiceErrNotification(err)
	}
	return &gostudyv1.NotificationResponse{Notification: notificationToProto(got)}, nil
}

func (s *NotificationsServer) Update(ctx context.Context, req *gostudyv1.UpdateNotificationRequest) (*gostudyv1.NotificationResponse, error) {
	upd := notificationFromProto(req.GetNotification())

	updated, err := service.UpdateNotification(int(req.GetId()), upd)
	if err != nil {
		return nil, mapServiceErrNotification(err)
	}

	return &gostudyv1.NotificationResponse{Notification: notificationToProto(updated)}, nil
}

func (s *NotificationsServer) Delete(ctx context.Context, req *gostudyv1.IdRequest) (*gostudyv1.DeleteResponse, error) {
	err := service.DeleteNotification(int(req.GetId()))
	if err != nil {
		return nil, mapServiceErrNotification(err)
	}
	return &gostudyv1.DeleteResponse{Deleted: true, Id: req.GetId()}, nil
}

func (s *NotificationsServer) List(ctx context.Context, req *gostudyv1.ListRequest) (*gostudyv1.ListNotificationsResponse, error) {
	notes := service.ListNotifications()
	out := make([]*gostudyv1.Notification, 0, len(notes))
	for _, n := range notes {
		out = append(out, notificationToProto(n))
	}
	return &gostudyv1.ListNotificationsResponse{Notifications: out}, nil
}
