package cache

type Cache[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (value V, ok bool)
	Len() int
	Cap() int
}
