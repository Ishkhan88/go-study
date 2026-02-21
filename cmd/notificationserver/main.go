package main

import (
	"log"
	"net"

	gostudyv1 "github.com/Ishkhan88/go-study/gen/gostudy/v1"
	transportgrpc "github.com/Ishkhan88/go-study/internal/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	s := grpc.NewServer()
	gostudyv1.RegisterNotificationsServiceServer(s, &transportgrpc.NotificationsServer{})

	log.Println("Notification gRPC server started on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
