package model

// структуры и интерфейсы
// Pet — это интерфейс "домашнее животное"
// Любой тип, у которого есть метод GetName() string будет считаться "Pet"

type Pet interface {
	GetName() string
}

type Cat struct {
	Name string
	Age  int
}

func (c Cat) GetName() string {
	return c.Name
}

type Dog struct {
	Name string
	Age  int
}

func (d Dog) GetName() string {
	return d.Name
}
