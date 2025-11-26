package main

import (
	"fmt"

	"github.com/Ishkhan88/go-study/internal/repository"
	"github.com/Ishkhan88/go-study/internal/service"
)

func main() {
	service.Run()

	fmt.Println("=== Все сохранённые коты ===")
	for _, c := range repository.GetCats() {
		fmt.Println(c.Name)
	}

	fmt.Println("=== Все сохранённые собаки ===")
	for _, d := range repository.GetDogs() {
		fmt.Println(d.Name)
	}
}
