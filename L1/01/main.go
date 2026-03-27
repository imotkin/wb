// Задача L1.1 - Встраивание структур
//
// Дана структура Human (с произвольным набором полей и методов).
// Реализовать встраивание методов в структуре Action
// от родительской структуры Human (аналог наследования).
//
// Подсказка: используйте композицию (embedded struct),
// чтобы Action имел все методы Human.

package main

import "fmt"

type Human struct {
	Name string
}

func (h Human) Hello() {
	fmt.Printf("Привет! Меня зовут %s\n", h.Name)
}

type Action struct {
	Human
}

func main() {
	a := Action{
		Human: Human{Name: "Илья"},
	}

	a.Hello()
}
