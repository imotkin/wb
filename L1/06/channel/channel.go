// Задача L1.6 - Остановка горутины
//
// В данном примере выполняется остановка горутины с помощью:
// (1) передачи значения пустой структуры в канал
// (2) закрытия переданного канала

package main

import (
	"log"
	"time"
)

// (1)
func exitWithChannelSend() {
	exit := make(chan struct{})

	go func() {
		for {
			select {
			case <-time.Tick(time.Second):
				log.Println("Function is active")
			case <-exit:
				return
			}
		}
	}()

	<-time.After(time.Millisecond * 3100)
	exit <- struct{}{}
}

// (2)
func exitWithChannelClose() {
	exit := make(chan struct{})

	go func() {
		for {
			select {
			case <-time.Tick(time.Second):
				log.Println("Function is active")
			case <-exit:
				return
			}
		}
	}()

	<-time.After(time.Millisecond * 3100)
	close(exit)
}

func main() {
	exitWithChannelClose()
}
