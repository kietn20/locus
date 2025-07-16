package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf(">>> GEOLOCATION EVENT: %s\n", msg.Payload())
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("event-listener-client") // A unique ID for this client

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %v", token.Error())
	}
	fmt.Println("Event listener connected to MQTT broker.")

	// This is the topic our geofence-service publishes events to.
	topic := "locus/geofence/events"
	// We use QoS 1 to ensure we get the messages even if there's a temporary network blip.
	if token := client.Subscribe(topic, 1, messagePubHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe: %v", token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)

	// Wait for a signal to gracefully disconnect.
	// This keeps the program running to listen for messages.
	fmt.Println("Waiting for events... Press Ctrl+C to exit.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Clean up on exit
	client.Unsubscribe(topic)
	client.Disconnect(250)
	fmt.Println("Event listener disconnected.")
}