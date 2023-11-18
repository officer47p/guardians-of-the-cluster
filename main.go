package main

import (
	"guardian/cache"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("failed to read env variables. Err: %s\n", err)
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		log.Panicf("failed to read env variable PORT. Err: %s\n", err)
	}
	defaultUserRequestQuota, err := strconv.ParseInt(os.Getenv("DEFAULT_USER_REQUEST_QUOTA"), 10, 64)
	if err != nil {
		log.Panicf("failed to read env variable DEFAULT_USER_REQUEST_QUOTA. Err: %s\n", err)
	}
	defaultUserTrafficQuota, err := strconv.ParseInt(os.Getenv("DEFAULT_USER_TRAFFIC_QUOTA"), 10, 64)
	if err != nil {
		log.Panicf("failed to read env variable DEFAULT_USER_TRAFFIC_QUOTA. Err: %s\n", err)
	}
	resetIntervalSeconds, err := strconv.ParseInt(os.Getenv("RESET_INTERVAL_SECONDS"), 10, 64)
	if err != nil {
		log.Panicf("failed to read env variable RESET_INTERVAL_SECONDS. Err: %s\n", err)
	}

	// Using in-memory cache, does not work when deploying multi instance of
	// this service

	// cache := cache.NewInMemoryCache()

	// Using redis-cache, does work when deploying multi instance of this service
	cache, err := cache.NewRedisCache()
	if err != nil {
		log.Panicf("failed to connect to redis. Err: %s\n", err)
	}

	rateLimiter := NewRateLimiter(
		cache,
		defaultUserRequestQuota,
		defaultUserTrafficQuota,
	)

	// Schedule cycle reset based on the provided reset interval
	ticker := time.NewTicker(time.Duration(resetIntervalSeconds * int64(time.Second)))
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err := rateLimiter.ResetCycle()
				if err != nil {
					log.Fatalf("error while resetting cycle. Err: %s", err)
				}
				log.Println("cycle was reset")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	httpServer := NewHttpServer(
		// HTTP server port
		int64(port),
		// Passing rate limiter middleware
		&rateLimiter,
		// Defining a dummy handler for success cases
		func(w http.ResponseWriter, r *http.Request) {
			_, err := io.WriteString(w, "Hello there :)")
			if err != nil {
				log.Printf("error handling the request. Err: %s\n", err)
			}
		},
	)

	log.Printf("server is listening on port %d\n", port)
	err = httpServer.Start()
	if err != nil {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
