package main

import (
	"log"
	"net"

	gostudyv1 "github.com/Ishkhan88/go-study/gen/gostudy/v1"
	transportgrpc "github.com/Ishkhan88/go-study/internal/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal("listen:", err)
	}

	s := grpc.NewServer()

	gostudyv1.RegisterUsersServiceServer(s, &transportgrpc.UsersServer{})
	gostudyv1.RegisterConcertsServiceServer(s, &transportgrpc.ConcertsServer{})
	gostudyv1.RegisterBookingsServiceServer(s, &transportgrpc.BookingsServer{})
	gostudyv1.RegisterNotificationsServiceServer(s, &transportgrpc.NotificationsServer{})

	log.Println("gRPC server started on :50053")
	if err := s.Serve(lis); err != nil {
		log.Fatal("serve:", err)
	}
}
