# About

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/asticode/go-ratelimiter/ratelimiter)

`go-ratelimiter` is a multi-buckets rate limiter for the GO programming language (http://golang.org).

# Install `go-ratelimiter`

Run the following command:

    $ go get github.com/asticode/go-ratelimiter/ratelimiter
    
# Example
    
    import (
        basememcache "github.com/bradfitz/gomemcache/memcache"
        "github.com/asticode/go-memcache/memcache"
        "github.com/asticode/go-ratelimiter/ratelimiter"
    )
    
    // Create the cache
    // In this example I'll use asticode/go-memcache but you can use another cache solution
    m := memcache.NewMemcache("myhost", "myprefix_", 10)
    
    // Create the rate limiter
    r := NewRateLimiter(
        m,
        basememcache.errCacheMiss,
    )
    
    // Add a limit of 2 requests for a duration of 5 seconds
    r.AddBucket(5, 2)
    
    // Add a limit of 10 request for a duration of 5 minutes
    r.AddBucket(5 * 60, 10)
    
    // Validate key rate limiter
    key := "key_test"
    e := r.Validate(key)
    e = r.Validate(key)
    e = r.Validate(key)
    
The last `validate` call will output an errLimitReached
