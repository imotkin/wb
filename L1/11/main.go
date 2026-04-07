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
	"cmp"
	"fmt"
	"slices"
)

func unique[T cmp.Ordered](nums ...T) map[T]struct{} {
	m := make(map[T]struct{})

	for _, n := range nums {
		m[n] = struct{}{}
	}

	return m
}

func intersection[T cmp.Ordered](a, b []T) []T {
	var (
		uniqueA = unique(a...)
		uniqueB = unique(b...)
		common  []T
	)

	for v := range uniqueA {
		if _, ok := uniqueB[v]; ok {
			common = append(common, v)
		}
	}

	slices.Sort(common)

	return common
}

func Print[T cmp.Ordered](a, b []T) {
	fmt.Printf("%v ∩ %v = %v\n", a, b, intersection[T](a, b))

}

func main() {
	intA := []int{1, 2, 3}
	intB := []int{2, 3, 4, 5}

	Print(intA, intB)

	strA := []string{"a", "b", "c", "d", "e"}
	strB := []string{"a", "c", "e"}

	Print(strA, strB)
}
