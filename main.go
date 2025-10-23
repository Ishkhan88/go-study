package main

import (
	"fmt"
)

// Функция для генерации шахматной доски
func chessBoard(size int) string {
	board := ""

	for y := 0; y < size; y++ { // строки
		for x := 0; x < size; x++ { // столбцы
			// Если сумма координат чётная — пробел, иначе решётка
			if (x+y)%2 == 0 {
				board += " "
			} else {
				board += "#"
			}
		}
		board += "\n" // переход на новую строку
	}

	return board
}

func main() {
	var size int
	fmt.Print("Введите размер доски (например 8): ")
	fmt.Scan(&size)

	result := chessBoard(size)
	fmt.Println(result)
}
