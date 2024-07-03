package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type HTTPServer struct {
	srv *http.Server
	mux *http.ServeMux
}

func NewHTTPServer(mux *http.ServeMux) (*HTTPServer, error) {
	mux.HandleFunc("/api/something", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Not sure what you expected to find here...")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}

	return &HTTPServer{
		srv: srv,
		mux: mux,
	}, nil
}

func (s *HTTPServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown: %s", err)
	}
	return nil
}

func (s *HTTPServer) Start() {
	// Start the HTTP server
	log.Println("Starting HTTP server on", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error running http server: %s\n", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		d, _ := httputil.DumpRequest(r, true)
		fmt.Println(string(d))

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log after the request is completed
		duration := time.Since(start)
		log.Printf("%s - %s %s - %s - %dms\n",
			r.RemoteAddr,
			r.Method,
			r.URL,
			r.UserAgent(),
			duration.Milliseconds(),
		)
	})
}
