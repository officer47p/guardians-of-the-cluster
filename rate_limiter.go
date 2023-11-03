package main

import "time"

// RateLimtier limits the number of requests per minute for each user based on
// some parameters of the request such as token and request size.
type RateLimiter struct {
	// Cache is the storage implementation that the RateLimiter service uses.
	Cache *Redis

	// DefaultUserRequestQuota is used when no quota limitiation is found for a
	// userId
	DefaultUserRequestQuota int64

	// DefaultUserTrafficQuota is used when no quota limitiation is found for a
	// userId
	DefaultUserTrafficQuota int64

	// Interval is the duration for a cycle. At the begining of each cycle,
	// rate limiter flushes the data in the cache, so rate-limited users can
	// continue sending requests again.
	Interval time.Duration
}

func NewRateLimiter(
	redis *Redis,
	defaultUserRequestQuota int64,
	defaultUserTrafficQuota int64,
	interval time.Duration,
) *RateLimiter {
	return &RateLimiter{Cache: redis,
		DefaultUserRequestQuota: defaultUserRequestQuota,
		DefaultUserTrafficQuota: defaultUserTrafficQuota,
		Interval:                interval,
	}
}
