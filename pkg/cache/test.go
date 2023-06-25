package cache

import (
	"testing"
)

func TestCache1(t *testing.T) {
	cache := new(Cache)
	cache.InitCache(10)
	for i := 0; i < 10; i++ {
		cache.Add(i)
		cache.Print()
	}
}
