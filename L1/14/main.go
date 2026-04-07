// L1.14 - Определение типа переменной в runtime
//
// Разработать программу, которая в runtime способна
// определить тип переменной, переданной в неё (на вход подаётся interface{}).
// Типы, которые нужно распознавать: int, string, bool, chan (канал).
// Подсказка: оператор типа switch v.(type) поможет в решении.

package main

import "fmt"

func defineSwitch(v any) string {
	switch v.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	default:
		return "chan"
	}
}

func defineSprint(v any) string {
	return fmt.Sprintf("%T", v)
}

func main() {
	values := []any{
		100,
		"101",
		true,
		make(chan int),
	}

	for _, v := range values {
		fmt.Printf(
			"%#v --> switch / %v, sprint / %v\n",
			v,
			defineSwitch(v),
			defineSprint(v),
		)
	}
}
