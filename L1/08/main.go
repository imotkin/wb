// Задача L1.8 - Установка бита в числе
//
// Дана переменная типа int64. Разработать программу, которая
// устанавливает i-й бит этого числа в 1 или 0.
//
// Пример: для числа 5 (0101₂) установка 1-го бита в 0 даст 4 (0100₂).
//
// Подсказка: используйте битовые операции (|, &, ^).

package main

import (
	"fmt"
)

func setBit(x, place, value int) int {
	if value == 0 {
		return x & ^(1 << (place - 1))
	}

	return x | (1 << (place - 1))
}

func main() {
	cases := []struct {
		x, place, value int
	}{
		{x: 0b1, place: 1, value: 0},   // 0b0
		{x: 0b100, place: 3, value: 0}, // 0b000
		{x: 0b100, place: 2, value: 1}, // 0b110
		{x: 0b100, place: 1, value: 1}, // 0b101
	}

	for _, c := range cases {
		fmt.Printf("%03b --> %03b\n", c.x, setBit(c.x, c.place, c.value))
	}
}
