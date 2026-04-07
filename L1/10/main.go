// L1.10 - Группировка температур
//
// Дана последовательность температурных колебаний:
// -25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5.
// Объединить эти значения в группы с шагом 10 градусов.
//
// Пример: -20:{-25.4, -27.0, -21.0}, 10:{13.0, 19.0, 15.5}, 20:{24.5}, 30:{32.5}.
//
// Пояснение: диапазон -20 включает значения от -20 до -29.9, диапазон 10 – от 10 до 19.9, и т.д.
// Порядок в подмножествах не важен.

package main

import (
	"fmt"
	"math"
	"strings"
)

func collect(temp ...float64) map[int][]float64 {
	m := make(map[int][]float64)

	var key int

	for _, t := range temp {
		if t < 0 {
			key = int(math.Ceil(t/10) * 10)
		} else {
			key = int(math.Floor(t/10) * 10)
		}

		m[key] = append(m[key], t)
	}

	return m
}

func main() {
	var (
		temp  = []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
		sets  = collect(temp...)
		parts = make([]string, 0, len(sets))
	)

	for group, temps := range sets {
		array := strings.Trim(fmt.Sprint(temps), "[]")
		parts = append(parts, fmt.Sprintf("%d:{%v}", group, array))
	}

	fmt.Println(strings.Join(parts, ", ") + ".")
}
