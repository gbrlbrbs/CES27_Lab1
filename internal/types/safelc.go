package types

import "sync"

type SafeLogicalClock struct {
	logicalClock int
	mu sync.Mutex
}

func (slc *SafeLogicalClock) Increase() {
	slc.mu.Lock()
	defer slc.mu.Unlock()

	slc.logicalClock++
}

func (slc *SafeLogicalClock) Set(val int) {
	slc.mu.Lock()
	defer slc.mu.Unlock()

	slc.logicalClock = val
}

func (slc *SafeLogicalClock) Get() int {
	return slc.logicalClock
}