// Задача L1.7 - Конкурентная запись в map
//
// Реализовать безопасную для конкуренции запись данных в структуру map.
//
// Подсказка: необходимость использования синхронизации
// (например, sync.Mutex или встроенная concurrent-map).
//
// Проверьте работу кода на гонки (util go run -race).

package main

import (
	"fmt"
	"sync"
)

type Map[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok = m.m[key]
	return
}

func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m[key] = value
}

func (m *Map[K, V]) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.m)
}

func main() {
	m := NewMap[int, string]()

	var wg sync.WaitGroup

	for i := range 1000 {
		wg.Go(func() {
			value := fmt.Sprintf("%d", i)
			m.Set(i, value)
		})
	}

	wg.Wait()

	fmt.Println(m.Len())
}
