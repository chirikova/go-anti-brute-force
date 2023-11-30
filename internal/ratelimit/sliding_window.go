package ratelimit

import (
	"sync"
	"time"
)

type window struct {
	limit          int64
	interval       time.Duration
	mu             sync.RWMutex
	timestamps     []int64
	lastAccessTime time.Time
}

func newWindow(limit int64, interval time.Duration) *window {
	return &window{
		limit:          limit,
		interval:       interval,
		timestamps:     make([]int64, 0, limit),
		lastAccessTime: time.Time{},
	}
}

func (w *window) add() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.timestamps = append(w.timestamps, time.Now().UnixNano())
	w.lastAccessTime = time.Now()
}

func (w *window) isLimitReached() bool {
	return w.size() >= w.limit
}

func (w *window) beginning() int {
	windowBeginning := time.Now().UnixNano() - w.interval.Nanoseconds()

	for i, value := range w.timestamps {
		if value >= windowBeginning {
			return i
		}
	}

	return 0
}

func (w *window) shift(index int) {
	w.timestamps = w.timestamps[index : len(w.timestamps)-1]
}

func (w *window) size() int64 {
	if beginning := w.beginning(); beginning > 0 {
		w.shift(beginning)
	}

	return int64(len(w.timestamps))
}
