// Задача L1.19 - Разворот строки
//
// Разработать программу, которая переворачивает подаваемую на вход строку.
// Например: при вводе строки «главрыба» вывод должен быть «абырвалг».
//
// Учтите, что символы могут быть в Unicode (русские буквы, emoji и пр.),
// то есть просто iterating по байтам может не подойти — нужен срез рун ([]rune).

package main

import (
	"fmt"
)

func main() {
	var s string
	fmt.Scan(&s)

	runes := []rune(s)

	for i := range runes {
		index := len(runes) - 1 - i
		fmt.Printf("%c", runes[index])
	}

	fmt.Println()
}
