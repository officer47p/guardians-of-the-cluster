package main

import (
	"errors"
	"fmt"
	"guardian/cache"
	"log"
)

// RateLimiter limits the number of requests per minute for each user based on
// some parameters of the request such as token and request size.
type RateLimiter struct {
	// Cache is the storage implementation that the RateLimiter service uses.
	cache cache.Cache

	// DefaultUserRequestQuota is used when no quota limitation is found for a
	// userId
	DefaultUserRequestQuota int64

	// DefaultUserTrafficQuota is used when no quota limitation is found for a
	// userId
	DefaultUserTrafficQuota int64
}

func NewRateLimiter(
	cache cache.Cache,
	defaultUserRequestQuota int64,
	defaultUserTrafficQuota int64,
) RateLimiter {
	return RateLimiter{cache: cache,
		DefaultUserRequestQuota: defaultUserRequestQuota,
		DefaultUserTrafficQuota: defaultUserTrafficQuota,
	}
}

func (rl *RateLimiter) ResetCycle() error {
	return rl.cache.FlushData()
}

func (rl *RateLimiter) CanMakeRequest(token string, dataSize int64) (bool, error) {
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

	// TODO: The code here can also be shorter, but I don't wanna spend much
	// time on it, maybe later :)
	// But for future fellow developer, you can put request number and data-
	// size related quota checks in two separate function and check for both
	// errors at once
	totalTrafficQuota, err := rl.getTotalTrafficQuota(token)
	if err != nil {
		log.Printf("error while calling getTotalTrafficQuota. Err: %s", err)
		return false, err
	}
	currentTrafficQuota, err := rl.getCurrentTrafficQuota(token)
	if err != nil {
		log.Printf("error while calling getCurrentTrafficQuota. Err: %s", err)
		return false, err
	}

	if currentTrafficQuota+dataSize > totalTrafficQuota {
		return false, nil
	}

	err = rl.setCurrentTrafficQuota(token, currentTrafficQuota+dataSize)
	if err != nil {
		log.Printf("error while setting the current traffic quota for token %s. Err: %s", token, err)
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
		if errors.Is(err, cache.ErrKeyNotFound) {
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
		if errors.Is(err, cache.ErrKeyNotFound) {
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

// TODO: Duplicated code for the below part, maybe refactor but keep it simple?
func (rl *RateLimiter) setCurrentTrafficQuota(token string, number int64) error {
	// quota:traffic:current:token
	return rl.cache.SetKey(fmt.Sprintf("quota:traffic:current:%s", token), number)

}

func (rl *RateLimiter) getTotalTrafficQuota(token string) (int64, error) {
	// quota:traffic:total:token
	q, err := rl.cache.GetKey(fmt.Sprintf("quota:traffic:total:%s", token))
	if err != nil {
		// If quota for user was not found
		if errors.Is(err, cache.ErrKeyNotFound) {
			// return default quota
			return rl.DefaultUserTrafficQuota, nil
		} else {
			log.Printf("error while getting the total traffic quota for token %s. Err: %s", token, err)
			return 0, err
		}

	}

	return q, nil
}

func (rl *RateLimiter) getCurrentTrafficQuota(token string) (int64, error) {
	// quota:traffic:current:token
	q, err := rl.cache.GetKey(fmt.Sprintf("quota:traffic:current:%s", token))
	if err != nil {
		// If quota for user was not found
		if errors.Is(err, cache.ErrKeyNotFound) {
			err := rl.cache.SetKey(fmt.Sprintf("quota:traffic:current:%s", token), 0)
			if err != nil {
				log.Printf("error setting the current traffic quota for token %s. Err: %s", token, err)
				return 0, err
			}

			return 0, nil
		}

		log.Printf("error while getting the current traffic quota for token %s. Err: %s", token, err)
		return 0, err
	}

	return q, nil
}
