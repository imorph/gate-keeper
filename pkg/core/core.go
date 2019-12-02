package core

import (
	"sync"
	"time"
)

// Limiter is counter with information on last update
type limiter struct {
	counter int       // current value of counter
	startTS time.Time // when counter was touched for the "first" time (1. first attempt ever 2. "bucket" refresh)
}

// LimitersCache is map of Limiter with mutex
type LimitersCache struct {
	threshold int           // threshold for counters (number of allowed tries per lifetime -- N M K)
	lifetime  time.Duration // lifetime of a Limiter/"bucket") usually minute
	hktime    time.Duration // House Keeping time interval (outdated entries gets whiped out once every hktime)
	mx        sync.RWMutex
	lm        map[string]*limiter
}

// NewCache creates new instance of Limiters cache
func NewCache(thr int, lt, hkt time.Duration) *LimitersCache {
	c := &LimitersCache{
		threshold: thr,
		lifetime:  lt,
		hktime:    hkt,
		lm:        make(map[string]*limiter),
	}
	return c
}

// Check will try to:
// 1. for given set of keys (login, pass, ip) lookup Limiters map
// 2. if key is not present create key and Limiter value (with counter=1 ant current timestamp)
// 3. if key is present check if current timestamp minus timestamp of last modification exeeds 1 minute
// 4. if yes resets counter
// 5. if no increments counter
// returns true if no brute-force detected returns false otherwise
func (l *LimitersCache) Check(key string) bool {
	currTS := time.Now()
	l.mx.Lock()
	val, ok := l.lm[key]
	defer l.mx.Unlock()
	if ok {
		if currTS.Sub(val.startTS) > l.lifetime {
			// it is more than a minute since this login/pass/ip was used last time
			// reseting counter to 1 and setting current ts
			l.lm[key] = &limiter{
				counter: 1,
				startTS: currTS,
			}
		} else {
			// it is at least second time this login pass or ip was used in this minute
			// incrementing counter and leave modification timestamp intact
			val.counter++
			l.lm[key] = &limiter{
				counter: val.counter,
				startTS: val.startTS,
			}
			if val.counter > l.threshold {
				// number of tries exceeded
				// updating startTS (this is new window)
				l.lm[key] = &limiter{
					counter: val.counter,
					startTS: val.startTS,
				}
				return false
			}
		}
	} else {
		// we did not found this login pass or ip in our map
		// so lets create it
		// first attempt and TS of window start
		l.lm[key] = &limiter{
			counter: 1,
			startTS: currTS,
		}
	}
	// since we are here we did not block anibody
	return true
}

// Reset will delete (if exist) counter for provided Key returns nothing since delete(map,key) also returns nothing
func (l *LimitersCache) Reset(key string) {
	l.mx.Lock()
	delete(l.lm, key)
	l.mx.Unlock()
}

// HouseKeep tries to get list of expired keys (first read-only map scan)
// and then delete all keys (that still outdated) from that list
func (l *LimitersCache) HouseKeep() {
	keysToDelete := l.mark()
	l.sweep(keysToDelete)
}

func (l *LimitersCache) mark() []string {
	currTS := time.Now()
	var keysToDelete []string
	l.mx.RLock()
	for key, value := range l.lm {
		if currTS.Sub(value.startTS) > l.hktime {
			keysToDelete = append(keysToDelete, key)
		}
	}
	l.mx.RUnlock()

	return keysToDelete
}

func (l *LimitersCache) sweep(keysToDelete []string) {
	currTS := time.Now()
	l.mx.Lock()
	for _, key := range keysToDelete {
		sts := l.lm[key].startTS
		if currTS.Sub(sts) > l.hktime {
			delete(l.lm, key)
		}
	}
	l.mx.Unlock()
}
