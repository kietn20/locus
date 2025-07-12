package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kietn20/locus/internal/vehicle"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var locationData vehicle.LocationData

	// Unmarshal the incoming message payload into our struct
	if err := locationData.FromJSON(msg.Payload()); err != nil {
		log.Printf("Error unmarshalling location data: %v", err)
		return
	}

	fmt.Printf("Received location for Vehicle '%s': Lat=%.4f, Lon=%.4f\n",
		locationData.VehicleID, locationData.Latitude, locationData.Longitude)
}

func main() {
	// --- MQTT Client Setup ---
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("location-service")
	opts.SetDefaultPublishHandler(messagePubHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %v", token.Error())
	}
	fmt.Println("Location service connected to broker.")

	// Subscribe to the wildcard topic to get data from all vehicles
	// The '+' is a single-level wildcard
	topic := "locus/vehicles/+/location"
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil { // Use QoS 1
		log.Fatalf("Failed to subscribe: %v", token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)

	// Wait for a signal to gracefully disconnect
	fmt.Println("Waiting for messages. Press Ctrl+C to exit.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Unsubscribe(topic)
	client.Disconnect(250)
	fmt.Println("Location service disconnected.")
}