package cache

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryCacheGetNonExistent(t *testing.T) {
	cache := New[string, int](10)

	_, ok := cache.Get("hello")

	require.False(t, ok)
}

func TestMemoryCacheGetExistent(t *testing.T) {
	cache := New[string, int](10)

	cache.Set("hello", 123)

	_, ok := cache.Get("hello")

	require.True(t, ok)
}

func TestMemoryCacheSet(t *testing.T) {
	cache := New[string, int](10)

	for i := range 10 {
		cache.Set(strconv.Itoa(i), i)
	}

	require.Equal(t, 10, cache.Len())

	for i := range 10 {
		v, ok := cache.Get(strconv.Itoa(i))
		require.True(t, ok)
		require.Equal(t, i, v)
	}
}

func TestMemoryCacheOverflow(t *testing.T) {
	cache := New[string, int](1)

	values := map[string]int{
		"hello": 123,
		"world": 456,
		"!":     789,
	}

	for k, v := range values {
		cache.Set(k, v)
	}

	require.Equal(t, 1, cache.Len())
	require.Equal(t, 1, cache.Cap())
}
