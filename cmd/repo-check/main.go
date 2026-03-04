package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Ishkhan88/go-study/internal/repository/postgres"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := postgres.New(ctx)
	if err != nil {
		panic(err)
	}
	defer db.Pool.Close()

	usersRepo := postgres.NewUsersRepository(db)

	u, err := usersRepo.CreateUser(
		ctx,
		"Test User",
		fmt.Sprintf("test_%d@mail.com", time.Now().Unix()),
	)
	if err != nil {
		panic(err)
	}

	got, err := usersRepo.GetUserByID(ctx, u.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created:", u.ID, u.Email, u.IsActive, u.Role)
	fmt.Println("Fetched:", got.ID, got.Email, got.IsActive, got.Role)
}
