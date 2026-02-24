package cache

import (
	"container/list"
	"sync"
)

type Entry struct {
	Key   any
	Value any
}

type MemoryCache[K comparable, V any] struct {
	mu       sync.RWMutex
	capacity int
	values   map[K]*list.Element
	queue    *list.List
}

func New[K comparable, V any](capacity int) *MemoryCache[K, V] {
	return &MemoryCache[K, V]{
		capacity: capacity,
		values:   make(map[K]*list.Element, capacity),
		queue:    list.New(),
	}
}

func (c *MemoryCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.values[key]; ok {
		c.queue.MoveToFront(element)
		element.Value.(*Entry).Value = value
		return
	}

	if c.queue.Len() == c.capacity {
		c.purge()
	}

	entry := &Entry{Key: key, Value: value}

	element := c.queue.PushFront(entry)
	c.values[key] = element
}

func (c *MemoryCache[K, V]) purge() {
	last := c.queue.Back()
	if last == nil {
		return
	}

	c.queue.Remove(last)

	delete(c.values, last.Value.(*Entry).Key.(K))
}

func (c *MemoryCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.values[key]
	if !ok {
		return value, ok
	}

	return v.Value.(*Entry).Value.(V), ok
}

func (c *MemoryCache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.queue.Len()
}

func (c *MemoryCache[K, V]) Cap() int {
	return c.capacity
}
