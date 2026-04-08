// Задача L1.26 - Уникальные символы в строке
//
// Разработать программу, которая проверяет, что все символы в строке встречаются один раз
// (т.е. строка состоит из уникальных символов).
// Вывод: true, если все символы уникальны, false, если есть повторения.
// Проверка должна быть регистронезависимой, т.е.
// символы в разных регистрах считать одинаковыми.
//
// Например: "abcd" -> true, "abCdefAaf" -> false (повторяются a/A), "aabcd" -> false.
//
// Подумайте, какой структурой данных удобно воспользоваться для проверки условия.

package main

import (
	"fmt"
	"unicode"
)

func unique(s string) bool {
	m := make(map[rune]struct{})

	for _, r := range s {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
		}

		if _, ok := m[r]; ok {
			return false
		}

		m[r] = struct{}{}
	}

	return true
}

func main() {
	var s string

	fmt.Scan(&s)

	fmt.Println(unique(s))
}
