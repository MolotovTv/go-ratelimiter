package ratelimiter

import (
	"testing"
	"time"

	"github.com/asticode/go-cache-manager/cachemanager"
	"github.com/stretchr/testify/assert"
)

func TestValidateCacheMiss(t *testing.T) {
	// Initialize
	k := "key_test"
	h := cachemanager.MockHandler()

	// Create rate limiter
	r := NewRateLimiter(h).AddBucket(time.Duration(2)*time.Second, 2)

	// Validate
	e := r.Validate(k)

	// Assert
	assert.NoError(t, e)
}

func TestValidateLimitReached(t *testing.T) {
	// Initialize
	k := "key_test"
	h := cachemanager.MockHandler()

	// Create rate limiter
	r := NewRateLimiter(h).AddBucket(time.Duration(2)*time.Second, 2)

	// Assert
	e := r.Validate(k)
	assert.NoError(t, e)
	e = r.Validate(k)
	assert.NoError(t, e)
	e = r.Validate(k)
	assert.EqualError(t, e, ErrLimitReached.Error())
}
