// Задача L1.20 - Разворот слов в предложении
//
// Разработать программу, которая переворачивает порядок слов в строке.
//
// Пример: входная строка:
// «snow dog sun», выход: «sun dog snow».
// Считайте, что слова разделяются одиночным пробелом. Постарайтесь
// не использовать дополнительные срезы, а выполнять операцию «на месте».

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()

	var (
		words = strings.Fields(s.Text())
		n     = len(words) - 1
	)

	for i := range n / 2 {
		words[i], words[n-i] = words[n-i], words[i]
	}

	fmt.Println(strings.Join(words, " "))
}
