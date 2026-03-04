package grpc

import gostudyv1 "github.com/Ishkhan88/go-study/gen/gostudy/v1"

type UsersServer struct {
	gostudyv1.UnimplementedUsersServiceServer
}

type ConcertsServer struct {
	gostudyv1.UnimplementedConcertsServiceServer
}

type BookingsServer struct {
	gostudyv1.UnimplementedBookingsServiceServer
	NotificationAddr string
}

type NotificationsServer struct {
	gostudyv1.UnimplementedNotificationsServiceServer
}
