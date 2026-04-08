// Задача L1.23 - Удаление элемента слайса
//
// Удалить i-ый элемент из слайса. Продемонстрируйте корректное удаление без утечки памяти.
//
// Подсказка: можно сдвинуть хвост слайса на место удаляемого элемента
// (copy(slice[i:], slice[i+1:])) и уменьшить длину слайса на 1.

package main

import (
	"fmt"
	"math/rand"
)

func Delete[T any](S []T, index int) []T {
	if index < 0 || index >= len(S) {
		return S
	}

	copy(S[index:], S[index+1:])
	return S[: len(S)-1 : len(S)-1]
}

func Print[T any](slice []T) {
	fmt.Printf("slice=%v len=%d cap=%d\n", slice, len(slice), cap(slice))
}

func main() {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	Print(slice)

	index := rand.Intn(len(slice))

	fmt.Printf("Delete -> a[%d] = %d\n", index, slice[index])

	slice = Delete(slice, index)
	Print(slice)
}
