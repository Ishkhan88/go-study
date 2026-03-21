// @title Auth Service API
// @version 1.0
// @description API для регистрации, логина и профиля пользователя
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/Ishkhan88/go-study/auth/docs"
	"github.com/Ishkhan88/go-study/auth/internal/config"
	"github.com/Ishkhan88/go-study/auth/internal/handler"
	"github.com/Ishkhan88/go-study/auth/internal/repository"
	"github.com/Ishkhan88/go-study/auth/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	cfg := config.Load()

	db, err := repository.NewPostgres(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	if err != nil {
		log.Fatal("database connection error:", err)
	}

	fmt.Println("Database connected")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Auth service is running")
	})

	mux.HandleFunc("/api/users", userHandler.Register)
	mux.HandleFunc("/api/auth/users", userHandler.Login)
	mux.HandleFunc("/api/users/profile", handler.AuthMiddleware(userHandler.Profile))

	fmt.Println("Auth service started on port 8081")

	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatal(err)
	}
}
