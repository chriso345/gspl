package main

import (
	"sync"
)

type handle uintptr

var (
	mu    sync.Mutex
	next  handle = 1
	store        = map[handle]any{}
)

func put(v any) handle {
	mu.Lock()
	defer mu.Unlock()
	h := next
	next++
	store[h] = v
	return h
}

func get(h handle) any {
	mu.Lock()
	defer mu.Unlock()
	return store[h]
}

func del(h handle) {
	mu.Lock()
	defer mu.Unlock()
	delete(store, h)
}
