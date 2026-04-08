// Задача L1.21 - Адаптер

package main

import "fmt"

type (
	Container interface {
		First() int
	}

	Array []int

	Adapter struct {
		Array
	}
)

func (a Array) Len() int {
	return len(a)
}

func (a Array) Cap() int {
	return len(a)
}

func (a Array) Get(i int) int {
	if i > a.Len() {
		return 0
	}

	return a[i]
}

func (a Adapter) First() int {
	return a.Array.Get(0)
}

func Print(c Container) {
	fmt.Printf("first -> %d\n", c.First())
}

func main() {
	a := Array([]int{1, 2, 3, 4})

	Print(Adapter{a})
}
