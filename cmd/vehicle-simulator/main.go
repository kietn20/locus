package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kietn20/locus/internal/vehicle"
)

func main() {
	// --- MQTT Client Setup ---
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		mqttBrokerHost = "localhost" // Fallback
	}
	mqttBrokerAddress := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	log.Printf("Connecting to MQTT broker at %s", mqttBrokerAddress)

	opts := mqtt.NewClientOptions().AddBroker(mqttBrokerAddress)
	// opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	
	// using the current timestamp to create a unique enough ID
	opts.SetClientID(fmt.Sprintf("vehicle-simulator-%d", time.Now().UnixNano()))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %v", token.Error())
	}
	fmt.Println("Vehicle simulator connected to broker.")

	// --- Simulation Logic ---
	vehicleID := "truck-01" // Let's simulate a specific truck.
	topic := fmt.Sprintf("locus/vehicles/%s/location", vehicleID)

	// Create a ticker that fires every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	// use a channel to listen for OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C: // This case executes every time the ticker fires
			// Create a new location data point with some random variation
			locationData := vehicle.LocationData{
				VehicleID: vehicleID,
				Latitude:  34.0522 + (rand.Float64()-0.5)*0.05, // Simulate movement
				Longitude: -118.2437 + (rand.Float64()-0.5)*0.05,
			}

			// Marshal the struct into JSON using our helper method.
			payload, err := locationData.ToJSON()
			if err != nil {
				log.Printf("Error converting to JSON: %v", err)
				continue // Skip this tick if there's an error
			}

			// Publish the JSON payload.
			token := client.Publish(topic, 1, false, payload) // Using QoS 1
			token.Wait()

			fmt.Printf("Published to %s: %s\n", topic, payload)

		case <-sigChan: // This case executes when we receive a shutdown signal
			fmt.Println("\nShutdown signal received. Disconnecting...")
			ticker.Stop()
			client.Disconnect(250)
			fmt.Println("Simulator disconnected.")
			return // Exit the program
		}
	}
}
