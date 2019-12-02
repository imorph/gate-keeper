package core

import (
	"fmt"
	"testing"
	"time"
)

func standartTry(cache *LimitersCache, wantRejects int, t *testing.T) {
	rejected := 0
	for i := 0; i < 10_000; i++ {
		if !cache.Check(fmt.Sprintf("key-%d", 1)) {
			rejected++
		}
	}
	if rejected != wantRejects {
		t.Errorf("rejected attempts=%d wantRejects=%d", rejected, wantRejects)
	}
}

func uniqueTry(cache *LimitersCache, wantRejects int, t *testing.T) {
	rejected := 0
	for i := 0; i < 10_000; i++ {
		if !cache.Check(fmt.Sprintf("key-%d", i)) {
			rejected++
		}
	}
	if rejected != wantRejects {
		t.Errorf("rejected attempts=%d wantRejects=%d", rejected, wantRejects)
	}
}

func TestNonUniqueChecksAllGood(t *testing.T) {
	cache := NewCache(10000, 1*time.Minute, 5*time.Minute)
	standartTry(cache, 0, t)
}

func TestBruteForceSlowpoke(t *testing.T) {
	cache := NewCache(10, 60*time.Millisecond, 5*time.Minute)
	rejected := 0
	wantRejects := 0
	for i := 0; i < 100; i++ {
		if !cache.Check(fmt.Sprintf("key-%d", 1)) {
			rejected++
		}
		time.Sleep(50 * time.Millisecond)
	}
	if rejected != wantRejects {
		t.Errorf("rejected attempts=%d wantRejects=%d", rejected, wantRejects)
	}
}

func TestEmptyHousekeep(t *testing.T) {
	cache := NewCache(10000, 1*time.Minute, 5*time.Minute)
	cache.HouseKeep()
}

func TestUniqueChecksAllGood(t *testing.T) {
	cache := NewCache(1, 1*time.Minute, 5*time.Minute)
	uniqueTry(cache, 0, t)
}

func TestRejectByAttempts(t *testing.T) {
	cache := NewCache(5000, 1*time.Minute, 5*time.Minute)
	standartTry(cache, 5000, t)
}

func TestBucketRefresh(t *testing.T) {
	cache := NewCache(10000, 100*time.Millisecond, 5*time.Minute)
	standartTry(cache, 0, t)

	time.Sleep(100 * time.Millisecond)
	// now possible to try another 10000 times
	standartTry(cache, 0, t)
}

func TestBucketHouseKeept(t *testing.T) {
	cache := NewCache(5000, 5*time.Minute, 1*time.Millisecond)
	standartTry(cache, 5000, t)

	time.Sleep(100 * time.Millisecond)
	// lets try another 10000 times
	// this ALL should fail because of 5 minute lifetime
	standartTry(cache, 10000, t)

	// lets expire our buckets
	time.Sleep(100 * time.Millisecond)
	// now housekeeping comes in
	cache.HouseKeep()
	// lets try another 10000 times
	// this 5000 should fail again because of threshold
	standartTry(cache, 5000, t)
}

func TestConcurrentUniqueTries(t *testing.T) {
	cache := NewCache(8, 5*time.Minute, 1*time.Millisecond)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	// wait to all settle
	time.Sleep(time.Second)
	// now all buckets should be "full" so no more new attemts allowed
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
}

func TestConcurrentBucketRefresh(t *testing.T) {
	cache := NewCache(8, 1*time.Second, 1*time.Millisecond)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	// wait to all settle
	time.Sleep(2 * time.Second)
	// now all buckets are refreshed and accept new tries
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
}

func TestConcurrentBucketRefreshWithFail(t *testing.T) {
	cache := NewCache(8, 1*time.Second, 1*time.Millisecond)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	// this all fail because buckets are full and not yet refreshed
	time.Sleep(800 * time.Millisecond)
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
	go uniqueTry(cache, 10000, t)
	// wait to all settle
	time.Sleep(2 * time.Second)
	// now all buckets are refreshed and accept new tries
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
	go uniqueTry(cache, 0, t)
}

func TestConcurrentChecksAndHouseKeeping(t *testing.T) {
	cache := NewCache(10000, 1*time.Second, 1*time.Millisecond)
	go func() {
		rejected := 0
		wantRejects := 0
		for i := 0; i < 8_000; i++ {
			if !cache.Check(fmt.Sprintf("key-%d", i)) {
				rejected++
			}
			time.Sleep(time.Millisecond)
		}
		if rejected != wantRejects {
			t.Errorf("rejected attempts=%d wantRejects=%d", rejected, wantRejects)
		}
	}()
	go func() {
		rejected := 0
		wantRejects := 0
		for i := 0; i < 8_000; i++ {
			if !cache.Check(fmt.Sprintf("key-%d", i)) {
				rejected++
			}
			time.Sleep(time.Millisecond)
		}
		if rejected != wantRejects {
			t.Errorf("rejected attempts=%d wantRejects=%d", rejected, wantRejects)
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			cache.HouseKeep()
		}
	}()
	time.Sleep(10 * time.Second)
}
