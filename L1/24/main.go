// Задача L1.24 - Расстояние между точками
//
// Разработать программу нахождения расстояния между двумя точками на плоскости.
// Точки представлены в виде структуры Point с инкапсулированными (приватными)
// полями x, y (типа float64) и конструктором. Расстояние рассчитывается
// по формуле между координатами двух точек.

// Подсказка: используйте функцию-конструктор NewPoint(x, y),
// Point и метод Distance(other Point) float64.

package main

import (
	"fmt"
	"math"
)

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

func (p Point) Distance(other Point) float64 {
	x := other.x - p.x
	y := other.y - p.y
	return math.Sqrt((x * x) + (y * y))
}

func (p Point) String() string {
	return fmt.Sprintf("(%.1f, %.1f)", p.x, p.y)
}

func main() {
	a := NewPoint(1, 10)
	b := NewPoint(2, 5)

	fmt.Printf("A: %s\n", a)
	fmt.Printf("B: %s\n", b)
	fmt.Printf("AB = %.0f\n", a.Distance(b))
}
