// Задача L1.6 - Остановка горутины
//
// В данном примере выполняется остановка горутины
// с помощью вызова функции Goexit из пакета runtime.

package main

import (
	"log"
	"runtime"
	"time"
)

func exitWithRuntime() {
	func() {
		after := time.After(time.Millisecond * 3100)
		for {
			select {
			case <-time.Tick(time.Second):
				log.Println("Function is active")
			case <-after:
				runtime.Goexit()
			}
		}
	}()
}

func main() {
	exitWithRuntime()
}
