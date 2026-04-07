package main

import "fmt"

func swapMath(a, b int) (int, int) {
	a = a + b
	b = a - b
	a = a - b

	return a, b
}

func swapXOR(a, b int) (int, int) {
	a = b ^ a
	b = a ^ b
	a = b ^ a

	return a, b
}

func main() {
	var a, b = 1, 100

	swap := []func(_, _ int) (int, int){
		swapMath, swapXOR,
	}

	for _, fn := range swap {
		sb, sa := fn(a, b)
		fmt.Printf("%d %d ⇆ %d %d\n", a, b, sb, sa)
	}
}
