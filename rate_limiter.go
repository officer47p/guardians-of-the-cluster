package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// RateLimtier limits the number of requests per minute for each user based on
// some parameters of the request such as token and request size.
type RateLimiter struct {
	// Cache is the storage implementation that the RateLimiter service uses.
	cache Cache

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
	cache Cache,
	defaultUserRequestQuota int64,
	defaultUserTrafficQuota int64,
	interval time.Duration,
) RateLimiter {
	return RateLimiter{cache: cache,
		DefaultUserRequestQuota: defaultUserRequestQuota,
		DefaultUserTrafficQuota: defaultUserTrafficQuota,
		Interval:                interval,
	}
}

func (rl *RateLimiter) ResetCycle() error {
	return rl.cache.FlushData()
}

func (rl *RateLimiter) CanMakeRequest(token string, _ int64) (bool, error) {
	totalRequestQuota, err := rl.getTotalRequestQuota(token)
	if err != nil {
		log.Printf("error while calling getTotalRequestQuota. Err: %s", err)
		return false, err
	}
	currentRequestQuota, err := rl.getCurrentRequestQuota(token)
	if err != nil {
		log.Printf("error while calling getCurrentRequestQuota. Err: %s", err)
		return false, err
	}

	if currentRequestQuota+1 > totalRequestQuota {
		return false, nil
	}

	err = rl.setCurrentRequestQuota(token, currentRequestQuota+1)
	if err != nil {
		log.Printf("error while setting the current request quota for token %s. Err: %s", token, err)
		return false, err
	}

	// Only passed here
	return true, nil

}

func (rl *RateLimiter) setCurrentRequestQuota(token string, number int64) error {
	// quota:request:current:token
	return rl.cache.SetKey(fmt.Sprintf("quota:request:current:%s", token), number)

}

func (rl *RateLimiter) getTotalRequestQuota(token string) (int64, error) {
	// quota:request:total:token
	q, err := rl.cache.GetKey(fmt.Sprintf("quota:request:total:%s", token))
	if err != nil {
		// If quota for user was not found
		if errors.Is(err, ErrKeyNotFound) {
			// return default quota
			return rl.DefaultUserRequestQuota, nil
		} else {
			log.Printf("error while getting the total request quota for token %s. Err: %s", token, err)
			return 0, err
		}

	}

	return q, nil
}

func (rl *RateLimiter) getCurrentRequestQuota(token string) (int64, error) {
	// quota:request:current:token
	q, err := rl.cache.GetKey(fmt.Sprintf("quota:request:current:%s", token))
	if err != nil {
		// If quota for user was not found
		if errors.Is(err, ErrKeyNotFound) {
			err := rl.cache.SetKey(fmt.Sprintf("quota:request:current:%s", token), 0)
			if err != nil {
				log.Printf("error setting the current request quota for token %s. Err: %s", token, err)
				return 0, err
			}

			return 0, nil
		}

		log.Printf("error while getting the current request quota for token %s. Err: %s", token, err)
		return 0, err
	}

	return q, nil
}
