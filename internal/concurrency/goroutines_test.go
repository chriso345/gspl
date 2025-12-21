package concurrency

import (
	"sync/atomic"
	"testing"
)

func TestTryAcquireAndRelease(t *testing.T) {
	atomic.StoreInt32(&MAX_GOROUTINES, 2)
	CURRENT_GOROUTINES.Store(0)

	if !TryAcquireGoroutine() {
		t.Fatal("expected first acquire to succeed")
	}
	if !TryAcquireGoroutine() {
		t.Fatal("expected second acquire to succeed")
	}
	if TryAcquireGoroutine() {
		t.Fatal("expected third acquire to fail")
	}

	ReleaseGoroutine()
	if !TryAcquireGoroutine() {
		t.Fatal("expected acquire to succeed after release")
	}

	CURRENT_GOROUTINES.Store(0)
}
