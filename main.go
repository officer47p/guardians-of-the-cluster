package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello there :)")
	if err != nil {
		log.Printf("error handling the request. Err: %s\n", err)
	}
}

func main() {
	port := 3333                           // TODO: Read from environment
	defaultUserRequestQuota := int64(5)    // TODO: Read from environment
	defaultUserTrafficQuota := int64(1024) // TODO: Read from environment

	cache := NewInMemoryCache()
	rateLimiter := NewRateLimiter(
		&cache,
		defaultUserRequestQuota,
		defaultUserTrafficQuota,
		time.Minute,
	)

	httpServer := NewHttpServer(int64(port), &rateLimiter, requestHandler)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go httpServer.Start()
	log.Printf("server is probably started on port %d\n", port)

	wg.Wait()
}
