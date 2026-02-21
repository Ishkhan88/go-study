package main

import (
	"log"
	"net"

	gostudyv1 "github.com/Ishkhan88/go-study/gen/gostudy/v1"
	transportgrpc "github.com/Ishkhan88/go-study/internal/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	s := grpc.NewServer()
	gostudyv1.RegisterBookingsServiceServer(s, &transportgrpc.BookingsServer{
		NotificationAddr: "127.0.0.1:50052",
	})

	log.Println("Booking gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
