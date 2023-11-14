package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := 3333                         // TODO: Read from environment
	defaultUserRequestQuota := int64(5)  // TODO: Read from environment
	defaultUserTrafficQuota := int64(10) // TODO: Read from environment
	resetInterval := time.Second * 10    // TODO: Read from environment

	cache := NewInMemoryCache()
	rateLimiter := NewRateLimiter(
		&cache,
		defaultUserRequestQuota,
		defaultUserTrafficQuota,
		resetInterval,
	)

	// Schedule cycle reset based on the provided reset interval
	ticker := time.NewTicker(resetInterval)
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
	err := httpServer.Start()
	if err != nil {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
