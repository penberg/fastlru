package lru

import (
	"testing"
)

func BenchmarkLRU(b *testing.B) {
	cache := NewCache()
	cache.Add("foo", "bar")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Get("foo")
		}
	})
}
