// Задача L1.22 - Большие числа и операции
//
// Разработать программу, которая перемножает, делит, складывает,
// вычитает две числовых переменных a, b, значения которых > 2^20 (больше 1 миллион).
//
// Комментарий: в Go тип int справится с такими числами,
// но обратите внимание на возможное переполнение для ещё больших значений.
// Для очень больших чисел можно использовать math/big.

package main

import (
	"fmt"
	"math"
	"math/big"
)

func add(a, b *big.Int) *big.Int {
	c := new(big.Int)
	return c.Add(a, b)
}

func sub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func div(a, b *big.Int) *big.Int {
	return new(big.Int).Div(a, b)
}

func mod(a, b *big.Int) *big.Int {
	return new(big.Int).Mod(a, b)
}

func main() {
	a := new(big.Int).SetUint64(math.MaxUint64)
	b := new(big.Int).SetUint64(math.MaxUint64)

	list := []struct {
		fn        func(a, b *big.Int) *big.Int
		operation string
	}{
		{
			fn:        add,
			operation: "+",
		},
		{
			fn:        sub,
			operation: "-",
		},
		{
			fn:        div,
			operation: "/",
		},
		{
			fn:        mod,
			operation: "%",
		},
	}

	for _, item := range list {
		c := item.fn(a, b)
		fmt.Printf("-> %d (a) %s %d (b) = %d\n", a, item.operation, b, c)
	}
}
