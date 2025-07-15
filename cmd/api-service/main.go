package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize the router and middleware 
	// chi is a lightweight, idiomatic and composable router for building Go HTTP services
	r := chi.NewRouter()
	r.Use(middleware.Logger) // for logging HTTP requests
	r.Use(middleware.Recoverer) // for recovering from panics and returning a 500 Internal Server Error
	
	// --- defining routes ---
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})
	})

	// --- Start the HTTP server ---
	port := "8080"
	log.Printf("API server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}