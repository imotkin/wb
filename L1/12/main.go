// L1.12 - Собственное множество строк
//
// Имеется последовательность строк: ("cat", "cat", "dog", "cat", "tree").
// Создать для неё собственное множество. Ожидается: получить набор уникальных слов.
// Для примера, множество = {"cat", "dog", "tree"}.

package main

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
)

type Set[T cmp.Ordered] map[T]struct{}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Values() []T {
	return slices.Collect(maps.Keys(s))
}

func unique[T cmp.Ordered](v ...T) Set[T] {
	m := make(Set[T])

	for _, t := range v {
		m[t] = struct{}{}
	}

	return m
}

func main() {
	words := []string{"cat", "cat", "dog", "cat", "tree"}

	set := unique(words...)

	fmt.Println(set.Values())
}
