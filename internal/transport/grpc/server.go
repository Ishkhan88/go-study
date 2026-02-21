package grpc

import gostudyv1 "github.com/Ishkhan88/go-study/gen/gostudy/v1"

// В Go нельзя одной структурой реализовать несколько gRPC-сервисов,
// если у них одинаковые имена методов (Create/Update/...), но разные типы запросов.
// Поэтому делаем отдельную структуру на каждый сервис.

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
