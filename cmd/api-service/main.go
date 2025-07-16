package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/kietn20/locus/internal/db"
	"github.com/kietn20/locus/internal/geofence"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	// "github.com/joho/godotenv"
)

type APIServer struct {
	Router *chi.Mux
	DBConn *pgx.Conn
}

func NewAPIServer(dbConn *pgx.Conn) *APIServer {
	s := &APIServer{
		Router: chi.NewRouter(),
		DBConn: dbConn,
	}

	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.setupRoutes()
	return s
}

func (s *APIServer) setupRoutes() {
	s.Router.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func (w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		r.Post("/geofences", s.createGeofenceHandler)
		r.Get("/geofences", s.listGeofencesHandler)
	})
}


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

// --- Handlers --- 

func (s *APIServer) createGeofenceHandler(w http.ResponseWriter, r *http.Request) {
	var feature geofence.GeoJSONFeature
	if err := json.NewDecoder(r.Body).Decode(&feature); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if feature.Properties.Name == "" || feature.Geometry.Type != "Polygon" {
		http.Error(w, "Missing name or invalid geometry type", http.StatusBadRequest)
		return
	}

	geomJSON, err := json.Marshal(feature.Geometry)
	if err != nil {
		http.Error(w, "Failed to process geometry", http.StatusInternalServerError)
		return
	}

	sql := `INSERT INTO geofences (name, area) VALUES ($1, ST_GeomFromGeoJSON($2))`
	_, err = s.DBConn.Exec(context.Background(), sql, feature.Properties.Name, string(geomJSON))
	if err != nil {
		log.Printf("Failed to insert geofence: %v", err)
		http.Error(w, "Failed to create geofence", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Geofence '%s' created successfully", feature.Properties.Name)
}

func (s *APIServer) listGeofencesHandler(w http.ResponseWriter, r *http.Request) {
	sql := `SELECT id, name, ST_AsGeoJSON(area), created_at FROM geofences ORDER BY created_at DESC`
	rows, err := s.DBConn.Query(context.Background(), sql)
	if err != nil {
		log.Printf("Failed to query geofences: %v", err)
		http.Error(w, "Failed to retrieve geofences", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var geofences []geofence.Geofence
	for rows.Next() {
		var gf geofence.Geofence
		if err := rows.Scan(&gf.ID, &gf.Name, &gf.Area, &gf.CreatedAt); err != nil {
			log.Printf("Failed to scan geofence row: %v", err)
			continue
		}
		geofences = append(geofences, gf)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(geofences)
}