package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	gostudyv1 "github.com/Ishkhan88/go-study/gen/gostudy/v1"
	"github.com/Ishkhan88/go-study/internal/config"
	"github.com/Ishkhan88/go-study/internal/repository"
	transportgrpc "github.com/Ishkhan88/go-study/internal/transport/grpc"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("env load warning:", err)
	}
	cfg := config.Default()

	// Mongo init
	if err := repository.InitMongo(context.Background(), cfg.MongoURI, cfg.MongoDB); err != nil {
		log.Fatal("mongo init error:", err)
	}
	defer func() {
		if err := repository.CloseMongo(context.Background()); err != nil {
			log.Println("mongo close error:", err)
		}
	}()

	// Redis init (audit)
	if err := repository.InitRedis(
		context.Background(),
		cfg.RedisAddr,
		cfg.RedisPassword,
		cfg.RedisDB,
		cfg.LogTTLSeconds,
	); err != nil {
		log.Fatal("redis init error:", err)
	}
	defer func() {
		if err := repository.CloseRedis(context.Background()); err != nil {
			log.Println("redis close error:", err)
		}
	}()

	// listener
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	s := grpc.NewServer()
	gostudyv1.RegisterNotificationsServiceServer(s, &transportgrpc.NotificationsServer{})

	// graceful stop
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Shutdown signal received... stopping Notification gRPC server")
		go func() {
			time.Sleep(5 * time.Second)
			s.Stop()
		}()
		s.GracefulStop()
	}()

	log.Println("Notification gRPC server started on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
