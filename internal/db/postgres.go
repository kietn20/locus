package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

// Connect initializes and returns a new database connection.
func Connect() *pgx.Conn {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var conn *pgx.Conn
	var err error

	// Retry logic
	for i := 0; i < 5; i++ {
		conn, err = pgx.Connect(context.Background(), dsn)
		if err == nil {
			// Success! Check if the connection is alive.
			err = conn.Ping(context.Background())
			if err == nil {
				log.Println("Successfully connected to database.")
				return conn
			}
		}

		log.Printf("Failed to connect to database (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second) //waiting for 2 seconds before retrying
	}

	// if we exit the loop, it means all retries failed
	log.Fatalf("Unable to connect to database after multiple retries: %v", err)
	return nil
}

// Migrate creates the necessary tables and indexes in the database.
func Migrate(conn *pgx.Conn) {
	// 'location' column of type 'GEOMETRY' from PostGIS extension
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS vehicle_locations (
		id SERIAL PRIMARY KEY,
		vehicle_id VARCHAR(255) NOT NULL,
		location GEOMETRY(Point, 4326) NOT NULL,
		timestamp TIMESTAMPTZ DEFAULT NOW() NOT NULL
	);

	-- Create an index on vehicle_id and timestamp for faster lookups.
	CREATE INDEX IF NOT EXISTS vehicle_locations_vehicle_id_timestamp_idx
	ON vehicle_locations (vehicle_id, timestamp DESC);
	
	-- Create a spatial index for fast location-based queries.
	CREATE INDEX IF NOT EXISTS vehicle_locations_location_idx
  	ON vehicle_locations USING GIST (location);
	`

	_, err := conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Fatalf("Table creation failed: %v\n", err)
	}


	// Create the geofences table with a spatial index for fast "contains" checks.
	// This table will store geofences defined by polygons.
	createGeofencesTableSQL := `
	CREATE TABLE IF NOT EXISTS geofences (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		area GEOMETRY(Polygon, 4326) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
	);
	-- Add a spatial index for fast "contains" checks later.
	CREATE INDEX IF NOT EXISTS geofences_area_idx ON geofences USING GIST (area);
	`
	_, err = conn.Exec(context.Background(), createGeofencesTableSQL)
	if err != nil {
		log.Fatalf("Geofences table creation failed: %v\n", err)
	}



	fmt.Println("Database migration completed successfully.")
}
