// Задача L1.6 - Остановка горутины
//
// В данном примере выполняется остановка горутины
// с помощью выхода из неё с return.

package main

import "log"

func exitWithReturn(flag bool) {
	log.Println("Function is active")

	if flag {
		return
	}
}

func main() {
	exitWithReturn(true)
}
