package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/kiet20/locus/internal/vehicle"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID(fmt.Sprintf("vehicle-simulator-%d", time.Now().UnixNano()))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %v", token.Error())
	}
	fmt.Println("Vehicle simulator connected to broker.")

	vehicleId := "truck-01"
	topic := fmt.Sprintf("locus/vehicles/%s/location", vehicleId)

	ticker := time.NewTicker(2 * time.Second)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			locationData := vehicle.LocationData{
				vehicleID: vehicleID,
				Latitude:  34.0522 + (rand.Float64()-0.5)*0.1, // Simulate movement
				Longitude: -118.2437 + (rand.float64()-0.5)*0.1,
			}

			payload, err := locationData.ToJSON()
			if err != nil {
				log.Printf("Error converting to JSON: %v", err)
				continue
			}

			token := client.Publish(topic, 1, false, payload)
			token.Wait()

			fmt.Printf("Published to %s: %s\n", topic, payload)

		case <-sigChan:
			fmt.Println("\nShutdown signal received. Disconnecting...")
			ticker.Stop()
			client.Disconnect(250)
			fmt.Println("Simulator disconnected.")
			return
		}
	}
}
