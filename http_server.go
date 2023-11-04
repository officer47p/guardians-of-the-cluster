package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type HttpServer struct {
	// These could also get into a different struct called HttpServerConfig,
	// but for sake of time, I just put them here.
	port        int64
	rateLimiter *RateLimiter
	mux         *http.ServeMux
}

func rateLimitMiddleware(next http.Handler, rateLimiter *RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("got new request")
		token := r.Host
		requestSize := r.ContentLength
		canMake, err := rateLimiter.CanMakeRequest(token, requestSize)
		if err != nil {
			log.Printf("error when validating request. Err: %s", err)
			io.WriteString(w, "internal-error")
			return
		}

		if !canMake {
			log.Println("request rate limited")
			io.WriteString(w, "rate-limited")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewHttpServer(port int64, rateLimiter *RateLimiter, handler http.HandlerFunc) *HttpServer {
	mux := http.NewServeMux()

	mux.Handle("/", rateLimitMiddleware(handler, rateLimiter))

	return &HttpServer{
		port:        port,
		rateLimiter: rateLimiter,
		mux:         mux,
	}
}

func (s *HttpServer) Start() {
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
	} else if err != nil {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
		// This could also be log.Fatalf, since it automatically calls os.exit(1)
	}
}
