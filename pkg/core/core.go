package core

import (
	"sync"
	"time"
)

// Limiter is counter with information on last update
type Limiter struct {
	counter uint16    // current value of counter
	lastTS  time.Time // when counter was touched last time (timestamp )
}

// LimitersCache is map of Limiter with mutex
type LimitersCache struct {
	threshold uint16        // threshold for counters (number of allowed tries per lifetime -- N M K)
	lifetime  time.Duration // lifetime of a Limiter/"bucket") usually minute
	hktime    time.Duration // House Keeping time interval (outdated entries gets whiped out once every hktime)
	mx        sync.RWMutex
	lm        map[string]Limiter
}

func NewCache(thr uint16, lt, hkt time.Duration) *LimitersCache {
	c := &LimitersCache{
		threshold: thr,
		lifetime:  lt,
		hktime:    hkt,
		lm:        make(map[string]Limiter),
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
	defer l.mx.Unlock()
	if val, ok := l.lm[key]; ok {
		if currTS.Sub(val.lastTS) > l.lifetime {
			// it is more than a minute since this login/pass/ip was used last time
			// reseting counter to 1 and setting current ts
			l.lm[key] = Limiter{
				counter: 1,
				lastTS:  currTS,
			}
		} else {
			// it is at least second time this login pass or ip was used in this minute
			// incrementing counter and set last modification to current TS
			val.counter++
			l.lm[key] = Limiter{
				counter: val.counter,
				lastTS:  currTS,
			}
			if val.counter > l.threshold {
				// number of tries exceeded
				return false
			}
		}
	} else {
		// we did not found this login pass or ip in our map
		// so lets create it
		l.lm[key] = Limiter{
			counter: 1,
			lastTS:  currTS,
		}
	}
	// since we are here we did not block anibody
	return true
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
		if currTS.Sub(value.lastTS) > l.hktime {
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
		if currTS.Sub(l.lm[key].lastTS) > l.hktime {
			delete(l.lm, key)
		}
	}
	l.mx.Unlock()
}
