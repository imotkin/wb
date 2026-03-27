// Задача L1.5 - Таймаут на канал
//
// Разработать программу, которая будет последовательно отправлять значения в канал,
// а с другой стороны канала – читать эти значения.
// По истечении N секунд программа должна завершаться.
//
// Подсказка: используйте time.After или таймер для ограничения времени работы.

package main

import (
	"context"
	"flag"
	"log"
	"time"
)

var (
	timeout  = flag.Duration("n", time.Second*5, "Exit timeout")
	interval = flag.Duration("interval", time.Second, "Message interval")
)

func runSender(ctx context.Context, to chan<- string, interval time.Duration) {
	for {
		select {
		case tick := <-time.Tick(interval):
			t := tick.Format("15:04:05")
			log.Println("-> Send tick:", t)
			to <- t
		case <-ctx.Done():
			return
		}
	}
}

func runReciever(ctx context.Context, from <-chan string) {
	for {
		select {
		case tick := <-from:
			log.Println("<-  Got tick:", tick)
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	ch := make(chan string)

	go runSender(ctx, ch, *interval)

	runReciever(ctx, ch)
}
