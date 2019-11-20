package core

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkCheckSingleThread(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Check(fmt.Sprintf("key-%d", i))
	}
}

func BenchmarkCheckParallel(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		b.ReportAllocs()

		for pb.Next() {
			cache.Check(fmt.Sprintf("key-%d", i))
			i++
		}
	})
}

func Benchmark0HKnoDeletes(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.HouseKeep()
	}
}

func Benchmark1000HKnoDeletes(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	for i := 0; i < 1000; i++ {
		cache.Check(fmt.Sprintf("key-%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.HouseKeep()
	}
}

func Benchmark10000HKnoDeletes(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	for i := 0; i < 10_000; i++ {
		cache.Check(fmt.Sprintf("key-%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.HouseKeep()
	}
}

func Benchmark100000HKnoDeletes(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	for i := 0; i < 100_000; i++ {
		cache.Check(fmt.Sprintf("key-%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.HouseKeep()
	}
}

func Benchmark1000000HKnoDeletes(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	for i := 0; i < 1_000_000; i++ {
		cache.Check(fmt.Sprintf("key-%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.HouseKeep()
	}
}

func Benchmark10000000HKnoDeletes(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	for i := 0; i < 10_000_000; i++ {
		cache.Check(fmt.Sprintf("key-%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.HouseKeep()
	}
}
