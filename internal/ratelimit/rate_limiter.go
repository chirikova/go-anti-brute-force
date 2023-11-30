package ratelimit

import (
	"time"
)

type RateLimiter struct {
	buckets  map[string]*window
	limit    int64
	interval time.Duration
	cancelCh <-chan struct{}
}

// NewRateLimiter
// The Sliding Window Algorithm is a rate-limiting technique
// that limits the number of requests a user can make within a given time frame
// while providing a smoother distribution of requests.
// It does this by continuously tracking requests
// and maintaining a "sliding window" that moves forward as time progresses,
// ensuring that request counts are always up-to-date.
func NewRateLimiter(interval time.Duration, limit int64, cancelCh <-chan struct{}) *RateLimiter {
	limiter := &RateLimiter{
		buckets:  make(map[string]*window),
		interval: interval,
		limit:    limit,
		cancelCh: cancelCh,
	}

	go func() {
		for {
			select {
			case <-time.After(interval):
				limiter.Clean()
			case <-cancelCh:
				return
			}
		}
	}()

	return limiter
}

func (r *RateLimiter) Allow(key string) bool {
	if _, ok := r.buckets[key]; !ok {
		r.buckets[key] = newWindow(r.limit, r.interval)
	}

	if !r.buckets[key].isLimitReached() {
		r.buckets[key].add()

		return true
	}

	return false
}

func (r *RateLimiter) Reset(key string) {
	delete(r.buckets, key)
}

func (r *RateLimiter) Clean() {
	for key, bucket := range r.buckets {
		if time.Since(bucket.lastAccessTime) > r.interval || bucket.size() == 0 {
			delete(r.buckets, key)
		}
	}
}
