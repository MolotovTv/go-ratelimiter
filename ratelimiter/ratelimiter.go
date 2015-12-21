package ratelimiter

import (
	"errors"
	"strconv"
)

var (
	errLimitReached = errors.New("Limit reached")
)

// RateLimiter represents a rate limiter
type RateLimiter interface {
	AddBucket(duration int, limit int) RateLimiter
	DelBucket(duration int) RateLimiter
	Validate(key string) error
}

// Cache represents the cache holding values typically a Memcache
type Cache interface {
	Set(k string, v []byte, ttl int32) error
	Decrement(key string, delta uint64) (uint64, error)
}

type rateLimiter struct {
	buckets map[int]int
	cache   Cache
	errCacheMiss error
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cache Cache, errCacheMiss error) RateLimiter {
	return &rateLimiter{
		buckets: make(map[int]int),
		cache:   cache,
		errCacheMiss: errCacheMiss,
	}
}

func (r *rateLimiter) AddBucket(duration int, limit int) RateLimiter {
	r.buckets[duration] = limit
	return r
}

func (r *rateLimiter) DelBucket(duration int) RateLimiter {
	delete(r.buckets, duration)
	return r
}

func (r rateLimiter) Validate(key string) error {
	// Loop through buckets
	for duration, limit := range r.buckets {
		// Decrement
		v, e := r.cache.Decrement(r.transformKey(key, duration), 1)

		// Process error
		if e != nil {
			// Cache miss
			if e.Error() == r.errCacheMiss.Error() {
				// Create key
				r.cache.Set(r.transformKey(key, duration), []byte(strconv.Itoa(limit)), int32(duration))
			} else {
				// Return
				return e
			}
		} else if v == 0 {
			return errLimitReached
		}
	}

	// Default return
	return nil
}

func (r rateLimiter) transformKey(key string, duration int) string {
	return key + "_" + strconv.Itoa(duration)
}
