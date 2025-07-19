package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kietn20/locus/internal/vehicle"
)

const (
	numVehicles = 50
)

func simulateVehicle(vehicleID string, client mqtt.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	topic := fmt.Sprintf("locus/vehicles/%s/location", vehicleID)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		// Create a new location data point with random movement.
		locationData := vehicle.LocationData{
			VehicleID: vehicleID,
			Latitude:  34.0522 + (rand.Float64()-0.5)*0.05,
			Longitude: -118.2437 + (rand.Float64()-0.5)*0.05,
		}

		payload, err := locationData.ToJSON()
		if err != nil {
			log.Printf("[%s] Error converting to JSON: %v", vehicleID, err)
			continue
		}

		// Publish the JSON payload.
		token := client.Publish(topic, 1, false, payload)
		// We don't wait for the token here to allow the simulator to run faster
		// and not get blocked by network latency.
		token.Wait()

		log.Printf("[%s] Published to %s: %s\n", vehicleID, topic, payload)

		// A small random sleep to make the vehicle movements less uniform.
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		<-ticker.C
	}
}

func main() {
	// --- MQTT Client Setup ---
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		mqttBrokerHost = "localhost"
	}
	mqttBrokerAddress := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	log.Printf("Connecting to MQTT broker at %s", mqttBrokerAddress)

	opts := mqtt.NewClientOptions().AddBroker(mqttBrokerAddress)
	opts.SetClientID("multi-vehicle-simulator") // unique client ID for the simulator
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %v", token.Error())
	}
	log.Println("Vehicle simulator connected to broker.")

	// --- Simulation Logic ---
	// the WaitGroup is used to wait for all goroutines to finish
	var wg sync.WaitGroup

	// launching a goroutine for each vehicle
	for i := 1; i <= numVehicles; i++ {
		vehicleID := fmt.Sprintf("truck-%02d", i)
		wg.Add(1) // increment the WaitGroup counter
		go simulateVehicle(vehicleID, client, &wg)
	}

	log.Printf("Launched %d vehicle simulators.", numVehicles)

	// waiting for a shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown signal received. Disconnecting...")
	client.Disconnect(250)
	log.Println("Simulator disconnected.")
}
