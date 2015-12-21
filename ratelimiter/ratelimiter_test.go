package ratelimiter

import (
	"math"
	"testing"
	"time"
	"strconv"

	"github.com/stretchr/testify/assert"
	"errors"
)

var (
	mockedErrCacheMiss = errors.New("Mocked cache miss")
)

type mockedCache struct {
	data       map[string][]byte
	expiration map[string]int32
}

func (c *mockedCache) Set(k string, v []byte, ttl int32) error {
	c.data[k] = v
	c.expiration[k] = int32(time.Now().Unix()) + ttl
	return nil
}

func (c *mockedCache) Decrement(key string, delta uint64) (uint64, error) {
	if int64(c.expiration[key]) < time.Now().Unix() {
		return 0, mockedErrCacheMiss
	} else {
		if _, ok := c.data[key]; !ok {
			return 0, mockedErrCacheMiss
		} else {
			v, _ := strconv.Atoi(string(c.data[key]))
			v -= 1
			c.data[key] = []byte(strconv.Itoa(v))
			return uint64(math.Max(float64(v), 0)), nil
		}
	}
}

func MockCache() *mockedCache {
	return &mockedCache{
		data:       make(map[string][]byte),
		expiration: make(map[string]int32),
	}
}

func TestValidateCacheMiss(t *testing.T) {
	// Initialize
	k := "key_test"
	c := MockCache()

	// Create rate limiter
	r := NewRateLimiter(
		c,
		mockedErrCacheMiss,
	).AddBucket(2, 2)

	// Validate
	e := r.Validate(k)

	// Assert
	assert.NoError(t, e)
	assert.Len(t, c.data, 1)
}

func TestValidateLimitReached(t *testing.T) {
	// Initialize
	k := "key_test"
	c := MockCache()

	// Create rate limiter
	r := NewRateLimiter(
		c,
		mockedErrCacheMiss,
	).AddBucket(2, 2)

	// Assert
	e := r.Validate(k)
	assert.NoError(t, e)
	e = r.Validate(k)
	assert.NoError(t, e)
	e = r.Validate(k)
	assert.Error(t, e)
}
