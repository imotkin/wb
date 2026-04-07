// Задача L1.6 - Остановка горутины
//
// В данном примере выполняется остановка горутины
// с помощью отмены после тайм-аута переданного контекста.

package main

import (
	"context"
	"log"
	"time"
)

func exitWithContext() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*3100)
	defer cancel()

	go func(ctx context.Context) {
		for {
			select {
			case <-time.Tick(time.Second):
				log.Println("Function is active")
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	<-ctx.Done()
}

func main() {
	exitWithContext()
}
