// Задача L1.17 - Бинарный поиск
//
// Реализовать алгоритм бинарного поиска встроенными методами языка.
// Функция должна принимать отсортированный слайс и искомый элемент,
// возвращать индекс элемента или -1, если элемент не найден.
//
// Подсказка: можно реализовать рекурсивно или итеративно, используя цикл for.

package main

import (
	"cmp"
	"fmt"
)

func binarySearch[T cmp.Ordered](slice []T, search T) int {
	var (
		left  = 0
		right = len(slice) - 1
	)

	for left <= right {
		idx := (left + right) / 2
		middle := slice[idx]

		switch {
		case search == middle:
			return idx
		case search > middle:
			left = idx + 1
		case search < middle:
			right = idx - 1
		}
	}

	return -1
}

func main() {
	nums := []int{0, 2, 4, 6, 8, 10}

	fmt.Println(binarySearch(nums, 1))
}
