// Задача L1.16 - Быстрая сортировка
//
// Реализовать алгоритм быстрой сортировки массива встроенными средствами языка.
// Можно использовать рекурсию.
//
// Подсказка: напишите функцию quickSort([]int) []int которая сортирует срез
// целых чисел. Для выбора опорного элемента можно взять середину или первый элемент.

package main

import (
	"fmt"
	"slices"
)

func quickSort(nums []int) []int {
	if len(nums) <= 1 {
		return nums
	}

	var less, equal, more []int

	pivot := nums[len(nums)/2]

	for _, v := range nums {
		switch {
		case v < pivot:
			less = append(less, v)
		case v > pivot:
			more = append(more, v)
		default:
			equal = append(equal, v)
		}
	}

	return slices.Concat(quickSort(less), equal, quickSort(more))
}

func main() {
	nums := []int{10, 2, 3, 7, 5, 1, 20, 0}

	fmt.Println(quickSort(nums))
}
