package ratelimiter

import (
	"errors"
	"strconv"
	"time"

	"github.com/asticode/go-cache-manager/cachemanager"
)

var (
	ErrLimitReached = errors.New("Limit reached")
)

// RateLimiter represents a rate limiter
type RateLimiter interface {
	AddBucket(duration time.Duration, limit int) RateLimiter
	SetBuckets(buckets map[time.Duration]int) RateLimiter
	DelBucket(duration time.Duration) RateLimiter
	Validate(key string) error
}

type rateLimiter struct {
	buckets map[time.Duration]int
	handler cachemanager.Handler
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(h cachemanager.Handler) RateLimiter {
	return &rateLimiter{
		buckets: make(map[time.Duration]int),
		handler: h,
	}
}

func (r *rateLimiter) AddBucket(duration time.Duration, limit int) RateLimiter {
	r.buckets[duration] = limit
	return r
}

func (r *rateLimiter) SetBuckets(buckets map[time.Duration]int) RateLimiter {
	r.buckets = buckets
	return r
}

func (r *rateLimiter) DelBucket(duration time.Duration) RateLimiter {
	delete(r.buckets, duration)
	return r
}

func (r rateLimiter) Validate(key string) error {
	// Loop through buckets
	for duration, limit := range r.buckets {
		// Decrement
		v, e := r.handler.Decrement(r.transformKey(key, duration), 1)

		// Process error
		if e != nil {
			// Cache miss
			if e.Error() == cachemanager.ErrCacheMiss.Error() {
				// Create key
				r.handler.Set(r.transformKey(key, duration), uint64(limit), duration)
			} else {
				// Return
				return e
			}
		} else if v == 0 {
			return ErrLimitReached
		}
	}

	// Default return
	return nil
}

func (r rateLimiter) transformKey(key string, duration time.Duration) string {
	return "ratelimiter:" + key + ":" + strconv.Itoa(int(duration))
}
