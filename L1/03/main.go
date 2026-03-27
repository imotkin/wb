// Задача L1.3 – Работа нескольких воркеров
//
// Реализовать постоянную запись данных в канал (в главной горутине).
// Реализовать набор из N воркеров, которые читают данные из этого канала и выводят их в stdout.
// Программа должна принимать параметром количество воркеров и
// при старте создавать указанное число горутин-воркеров.

package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"
)

var (
	n        = flag.Int("n", 5, "Number of workers")
	interval = flag.Duration("interval", time.Second, "Worker function interval")
)

func runWorkers(ch <-chan string, n int) {
	for i := range n {
		go func() {
			for word := range ch {
				fmt.Printf("Worker %d got %q\n", i+1, word)
			}
		}()
	}
}

func runSender(ctx context.Context, to chan<- string, interval time.Duration, words []string) {
	var i int

	for {
		select {
		case <-time.Tick(interval):
			to <- words[i%4]
			i++
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	words := make(chan string)

	go runWorkers(words, *n)

	runSender(
		ctx,
		words,
		*interval,
		[]string{"Hello", ",", "World", "!"},
	)
}
