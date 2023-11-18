package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
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
		userId, ok := r.Header["X-User-Id"]
		if !ok || len(userId) == 0 {
			log.Println("error when reading user id from header. User ID does not exist")
			io.WriteString(w, "user-id-does-not-exist")
			return
		}
		// TODO: Implement correct token from headers.
		token := userId[0]
		requestSize := r.ContentLength

		log.Printf("Token: %s\n", token)
		log.Printf("Request Size: %d\n", requestSize)

		canMake, err := rateLimiter.CanMakeRequest(token, requestSize)
		// TODO: Return status code and proper response structure
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

func (s *HttpServer) Start() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
		return nil
	}

	return err
}
