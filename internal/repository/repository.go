package repository

import (
	"fmt"

	"github.com/Ishkhan88/go-study/internal/model"
)

//"хранение" данных
//объявление слайсов и функций, которые складывают котов и собак

var cats []model.Cat
var dogs []model.Dog

func SavePet(p model.Pet) {
	switch v := p.(type) {
	case model.Cat:
		cats = append(cats, v)
		fmt.Println("Добавили кота:", v.Name)

	case model.Dog:
		dogs = append(dogs, v)
		fmt.Println("Добавили собаку:", v.Name)

	default:
		fmt.Println("Неизвестный тип животного")
	}
}

func GetCats() []model.Cat { //получить и посмотреть всех котов
	return cats
}

func GetDogs() []model.Dog {
	return dogs
}
