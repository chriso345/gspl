package concurrency

import (
	"runtime"
	"sync/atomic"
)

// MAX_GOROUTINES defines the maximum number of parallel goroutines for branching.
// It defaults to number of CPUs but can be adjusted by the caller before solving.
var MAX_GOROUTINES int32

// CURRENT_GOROUTINES tracks the current number of active branching goroutines.
var CURRENT_GOROUTINES atomic.Int32

func init() {
	maxGoroutines := int32(getPhysicalCores())
	atomic.StoreInt32(&MAX_GOROUTINES, maxGoroutines)
}

// getPhysicalCores returns the number of physical CPU cores available.
// This is a rough estimate, as most systems have 2 threads per core.
func getPhysicalCores() int {
	return runtime.NumCPU() / 2
}

// TryAcquireGoroutine attempts to increment CURRENT_GOROUTINES if below MAX_GOROUTINES.
// Returns true if a goroutine slot was acquired, false otherwise.
func TryAcquireGoroutine() bool {
	max := atomic.LoadInt32(&MAX_GOROUTINES)
	if max <= 0 {
		return false
	}
	for {
		cur := CURRENT_GOROUTINES.Load()
		if cur >= max {
			return false
		}
		if CURRENT_GOROUTINES.CompareAndSwap(cur, cur+1) {
			return true
		}
	}
}

// ReleaseGoroutine decrements the CURRENT_GOROUTINES counter.
func ReleaseGoroutine() {
	CURRENT_GOROUTINES.Add(-1)
}
