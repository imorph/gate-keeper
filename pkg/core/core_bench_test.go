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

func Benchmark0Mark(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cache.mark()
	}
}

func Benchmark1000Mark(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	for i := 0; i < 1000; i++ {
		cache.Check(fmt.Sprintf("key-%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cache.mark()
	}
}

func Benchmark10000Mark(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	for i := 0; i < 10_000; i++ {
		cache.Check(fmt.Sprintf("key-%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cache.mark()
	}
}

func Benchmark100000Mark(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	for i := 0; i < 100_000; i++ {
		cache.Check(fmt.Sprintf("key-%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cache.mark()
	}
}

func Benchmark0Sweep(b *testing.B) {
	cache := NewCache(100, 5*time.Minute, 5*time.Minute)
	zero := cache.mark()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.sweep(zero)
	}
}

func Benchmark1000Sweep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cache := NewCache(100, 5*time.Microsecond, 5*time.Microsecond)
		for i := 0; i < 1000; i++ {
			cache.Check(fmt.Sprintf("key-%d", i))
		}
		delete := cache.mark()
		b.StartTimer()
		cache.sweep(delete)
	}
}

// sweeps on my machine ~5k-6k records
func Benchmark10000Sweep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cache := NewCache(100, 5*time.Microsecond, 3000*time.Microsecond)
		for i := 0; i < 10_000; i++ {
			cache.Check(fmt.Sprintf("key-%d", i))
			if i == 5000 {
				time.Sleep(time.Microsecond)
			}

		}
		delete := cache.mark()
		b.StartTimer()
		cache.sweep(delete)
	}
}
