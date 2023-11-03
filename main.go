package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func requestHanlder(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "Hello there :)")
	if err != nil {
		log.Printf("error handling the request. Err: %s\n", err)
	}
}

func main() {
	wg := sync.WaitGroup{}
	port := 3333 // TODO: Read from environment
	redis := NewRedis()
	rateLimiter := NewRateLimiter(redis, 100, 1024, time.Minute)

	httpServer := NewHttpServer(int64(port), rateLimiter, requestHanlder)

	wg.Add(1)
	go httpServer.Start()
	log.Printf("server is probably started on port %d\n", port)

	wg.Wait()
}
