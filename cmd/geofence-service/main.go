package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/kietn20/locus/internal/db"
	"github.com/kietn20/locus/internal/vehicle"
)

type GeofenceService struct {
	DBConn        *pgx.Conn
	MQTTClient    mqtt.Client
	vehicleStates map[string]string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// --- Connect to Database ---
	conn := db.Connect()
	defer conn.Close(context.Background())

	service := &GeofenceService{
		DBConn:        conn,
		vehicleStates: make(map[string]string),
	}

	// --- Connect to MQTT ---
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("geofence-service")
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	fmt.Println("Geofence service connected to MQTT broker.")

	service.MQTTClient = mqttClient

	// opts.SetDefaultPublishHandler(service.messageHandler)

	topic := "locus/vehicles/+/location"
	if token := mqttClient.Subscribe(topic, 1, service.messageHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe: %v", token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)

	fmt.Println("Waiting for messages. Press Ctrl+C to exit.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mqttClient.Unsubscribe(topic)
	mqttClient.Disconnect(250)
	fmt.Println("Geofence service disconnected.")
}

// This function performs the core geospatial query.
func (s *GeofenceService) checkGeofence(loc vehicle.LocationData) (string, error) {
	var geofenceName string

	// ST_Contains checks if the first geometry contains the second.
	// We check if any geofence's 'area' contains the vehicle's 'location' point.
	// We return the name of the first one we find. LIMIT 1 makes it efficient.
	sql := `
		SELECT name FROM geofences 
		WHERE ST_Contains(area, ST_SetSRID(ST_MakePoint($1, $2), 4326))
		LIMIT 1
	`
	// Note: We use longitude first, then latitude for PostGIS point creation.
	err := s.DBConn.QueryRow(context.Background(), sql, loc.Longitude, loc.Latitude).Scan(&geofenceName)

	if err == pgx.ErrNoRows {
		// This is not a system error. It's the expected result when the point is outside all geofences.
		return "", nil
	}
	if err != nil {
		// This is a real database error.
		return "", err
	}

	return geofenceName, nil
}

// messageHandler is the callback for processing vehicle location updates.
func (s *GeofenceService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	var loc vehicle.LocationData
	if err := loc.FromJSON(msg.Payload()); err != nil {
		log.Printf("Error unmarshalling location data: %v", err)
		return
	}

	// Get the previous state of the vehicle from our in-memory map.
	previousFence, _ := s.vehicleStates[loc.VehicleID]

	// Check which geofence the vehicle is currently in.
	currentFence, err := s.checkGeofence(loc)
	if err != nil {
		log.Printf("Error checking geofence for vehicle %s: %v", loc.VehicleID, err)
		return
	}

	// Compare the previous state with the current state.
	if currentFence != previousFence {
		// The state has changed!
		if previousFence != "" {
			// The vehicle has exited a fence.
			eventPayload := fmt.Sprintf(`{"vehicle_id": "%s", "geofence_name": "%s", "event": "exit"}`, loc.VehicleID, previousFence)
			fmt.Printf("EVENT: Vehicle %s exited geofence %s\n", loc.VehicleID, previousFence)
			s.MQTTClient.Publish("locus/geofence/events", 1, false, eventPayload)
		}
		if currentFence != "" {
			// The vehicle has entered a fence.
			eventPayload := fmt.Sprintf(`{"vehicle_id": "%s", "geofence_name": "%s", "event": "enter"}`, loc.VehicleID, currentFence)
			fmt.Printf("EVENT: Vehicle %s entered geofence %s\n", loc.VehicleID, currentFence)
			s.MQTTClient.Publish("locus/geofence/events", 1, false, eventPayload)
		}

		// Update the state map with the new location.
		s.vehicleStates[loc.VehicleID] = currentFence
	}
}
