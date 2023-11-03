package main

import (
	"errors"
	"fmt"
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

func NewHttpServer(port int64, rateLimiter *RateLimiter, handler http.HandlerFunc) *HttpServer {
	mux := http.NewServeMux()

	mux.Handle("/", handler)

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
