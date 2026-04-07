// Задача L1.18 - Конкурентный счетчик
//
// Реализовать структуру-счётчик, которая будет инкрементироваться
// в конкурентной среде (т.е. из нескольких горутин). По завершению
// программы структура должна выводить итоговое значение счётчика.
//
// Подсказка: вам понадобится механизм синхронизации, например,
// sync.Mutex или sync/Atomic для безопасного инкремента.

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Counter struct {
	i atomic.Int64
}

func (c *Counter) Inc() {
	c.i.Add(1)
}

func (c *Counter) Value() int64 {
	return c.i.Load()
}

func main() {
	var (
		c  Counter
		wg sync.WaitGroup
	)

	for range 1000 {
		wg.Go(func() {
			c.Inc()
		})
	}

	wg.Wait()

	fmt.Println(c.Value())
}
