// L1.11 - Пересечение множеств
//
// Реализовать пересечение двух неупорядоченных множеств (например, двух слайсов) —
// т.е. вывести элементы, присутствующие и в первом, и во втором.
//
// Пример:
// A = {1,2,3}
// B = {2,3,4}
// Пересечение = {2,3}

package main

import (
	"fmt"
	"slices"
)

func unique(nums ...int) map[int]struct{} {
	m := make(map[int]struct{})

	for _, n := range nums {
		m[n] = struct{}{}
	}

	return m
}

func main() {
	a := []int{1, 2, 3}
	b := []int{2, 3, 4, 5}

	uniqueA := unique(a...)
	uniqueB := unique(b...)

	var common []int

	for num := range uniqueA {
		if _, ok := uniqueB[num]; ok {
			common = append(common, num)
		}
	}

	slices.Sort(common)

	fmt.Println(common)
}
