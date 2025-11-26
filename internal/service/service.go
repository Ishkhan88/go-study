package service

import (
	"fmt"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"github.com/Ishkhan88/go-study/internal/repository"
)

func Run() {
	for i := 1; i <= 3; i++ {
		fmt.Println("---- Итерация", i, "----")

		cat := model.Cat{
			Name: fmt.Sprintf("Снежок-%d", i),
		}

		dog := model.Dog{
			Name: fmt.Sprintf("Шарик-%d", i),
		}

		repository.SavePet(cat)
		repository.SavePet(dog)

		time.Sleep(1 * time.Second)
	}
}
