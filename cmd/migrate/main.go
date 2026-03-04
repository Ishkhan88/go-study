package main

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL is empty")
	}

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No changes (already migrated).")
			return
		}
		panic(err)
	}

	fmt.Println("Migrations applied successfully.")
}
