package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

// Connect initializes and returns a new database connection.
func Connect() *pgx.Conn {
	// Construct the Database Source Name (DSN) string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Pinging the database to verify the connection to ensure it's alive
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	fmt.Println("Successfully connected to database.")
	return conn
}

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

	fmt.Println("Database migration completed successfully.")
}
