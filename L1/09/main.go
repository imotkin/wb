// Задача L1.9 - Конвейер чисел
//
// Разработать конвейер чисел. Даны два канала:
// в первый пишутся числа x из массива, во второй – результат операции x*2.
// После этого данные из второго канала должны выводиться в stdout.
// То есть, организуйте конвейер из двух этапов с горутинами:
// генерация чисел и их обработка. Убедитесь, что чтение из второго канала корректно завершается.

package main

import (
	"flag"
	"fmt"
)

func numbers(nums []int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _, v := range nums {
			out <- v
		}
	}()

	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for num := range in {
			out <- num * num
		}
	}()

	return out
}

func output(in <-chan int) {
	for num := range in {
		fmt.Println(num)
	}
}

var n = flag.Int("n", 5, "Number of array elements")

func main() {
	flag.Parse()

	array := make([]int, *n)

	for i := range *n {
		array[i] = i + 1
	}

	nums := numbers(array)
	squared := square(nums)
	output(squared)
}
