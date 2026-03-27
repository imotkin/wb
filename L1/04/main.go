// Задача L1.4 - Завершение по Ctrl+C
//
// Программа должна корректно завершаться по нажатию Ctrl+C (SIGINT).
// Выберите и обоснуйте способ завершения работы всех горутин-воркеров при получении сигнала прерывания.
//
// Подсказка: можно использовать контекст (context.Context) или канал для оповещения о завершении.

package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	n        = flag.Int("n", 5, "Number of workers")
	interval = flag.Duration("interval", time.Second, "Worker function interval")
	timeout  = flag.Duration("timeout", time.Second*5, "Exit timeout")
)

func runWorker(ctx context.Context, id int, interval time.Duration) {
	start := time.Now()

	for {
		select {
		case <-time.Tick(interval):
			log.Printf("Worker %d is active...\n", id)
		case <-ctx.Done():
			log.Printf("Worker %d was stopped after %.0fs\n", id, time.Since(start).Seconds())
			return
		}
	}
}

func run(ctx context.Context, n int, interval, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var wg sync.WaitGroup

	for i := 1; i <= n; i++ {
		wg.Go(func() { runWorker(ctx, i, interval) })
	}

	wg.Wait()
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	run(ctx, *n, *interval, *timeout)
}
