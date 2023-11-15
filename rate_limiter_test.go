package main

import (
	"testing"
	"time"
)

var (
	defaultUserRequestQuota = int64(5)         // TODO: Read from environment
	defaultUserTrafficQuota = int64(10)        // TODO: Read from environment
	resetInterval           = time.Second * 10 // TODO: Read from environment
)

func TestRateLimiterRequestNumberQuota(t *testing.T) {

	t.Run("should rate-limit", func(t *testing.T) {
		token := "some-string-as-token"
		cache := NewInMemoryCache()
		rateLimiter := NewRateLimiter(
			&cache,
			defaultUserRequestQuota,
			defaultUserTrafficQuota,
			resetInterval,
		)

		for i := 0; i < int(defaultUserRequestQuota); i++ {
			ok, err := rateLimiter.CanMakeRequest(token, 1)
			if err != nil {
				t.Fatalf("hit error. Err: %s", err)
			}
			if !ok {
				t.Fatal("rate limited")
			}
		}

		ok, err := rateLimiter.CanMakeRequest(token, 10)
		if err != nil {
			t.Fatalf("hit error")
		}
		if ok {
			t.Fatalf("was not rate limited")
		}
	})

	t.Run("should recover from rate-limit after cycle reset", func(t *testing.T) {
		token := "some-string-as-token"
		cache := NewInMemoryCache()
		rateLimiter := NewRateLimiter(
			&cache,
			defaultUserRequestQuota,
			defaultUserTrafficQuota,
			resetInterval,
		)

		for i := 0; i < int(defaultUserRequestQuota*2); i++ {
			_, err := rateLimiter.CanMakeRequest(token, 0)
			if err != nil {
				t.Fatalf("hit error. Err: %s", err)
			}
		}

		ok, err := rateLimiter.CanMakeRequest(token, 0)
		if err != nil {
			t.Fatalf("hit error")
		}
		if ok {
			t.Fatalf("was not rate limited")
		}

		err = rateLimiter.ResetCycle()
		if err != nil {
			t.Fatalf("fail to reset cycle")
		}

		ok, err = rateLimiter.CanMakeRequest(token, 1)
		if err != nil {
			t.Fatalf("hit error")
		}
		if !ok {
			t.Fatalf("was rate limited")
		}
	})

}

func TestRateLimiterTrafficQuota(t *testing.T) {
	t.Run("should-rate-limit", func(t *testing.T) {
		token := "some-string-as-token"

		cache := NewInMemoryCache()
		rateLimiter := NewRateLimiter(
			&cache,
			defaultUserRequestQuota,
			defaultUserTrafficQuota,
			resetInterval,
		)

		for i := 0; i < 2; i++ {
			ok, err := rateLimiter.CanMakeRequest(token, 4)
			if err != nil {
				t.Fatalf("hit error. Err: %s", err)
			}
			if !ok {
				t.Fatal("rate limited")
			}
		}
		ok, err := rateLimiter.CanMakeRequest(token, 4)
		if err != nil {
			t.Fatalf("hit error")
		}
		if ok {
			t.Fatalf("was not rate limited")
		}
	})

	t.Run("should recover from rate-limit after cycle reset", func(t *testing.T) {
		token := "some-string-as-token"

		cache := NewInMemoryCache()
		rateLimiter := NewRateLimiter(
			&cache,
			defaultUserRequestQuota,
			defaultUserTrafficQuota,
			resetInterval,
		)

		for i := 0; i < 2; i++ {
			ok, err := rateLimiter.CanMakeRequest(token, defaultUserTrafficQuota/2)
			if err != nil {
				t.Fatalf("hit error. Err: %s", err)
			}
			if !ok {
				t.Fatal("rate limited")
			}
		}
		ok, err := rateLimiter.CanMakeRequest(token, defaultUserTrafficQuota/2)
		if err != nil {
			t.Fatalf("hit error")
		}
		if ok {
			t.Fatalf("was not rate limited")
		}

		err = rateLimiter.ResetCycle()
		if err != nil {
			t.Fatalf("fail to reset cycle")
		}

		ok, err = rateLimiter.CanMakeRequest(token, 1)
		if err != nil {
			t.Fatalf("hit error")
		}
		if !ok {
			t.Fatalf("was rate limited")
		}
	})

}
