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

func TestRateLimiterReturnErrorAfterHittingRequestNumberQuota(t *testing.T) {
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

}

func TestRateLimiterReturnErrorAfterHittingTrafficQuota(t *testing.T) {
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

}
