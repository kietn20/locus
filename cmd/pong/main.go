package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// messagePubHandler is a callback function that will be called when a message is received on the subscribed topic
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: '%s' from topic: '%s'\n", msg.Payload(), msg.Topic())
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("go-subscriber")
	// Set a callback that will be called when the client is reconnected.
	opts.SetDefaultPublishHandler(messagePubHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %v", token.Error())
	}
	fmt.Println("Subscriber connected to broker.")

	// Subscribe to the topic
	topic := "locus/test"
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe: %v", token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)

	fmt.Println("Waiting for messages. Press Ctrl+C to exit.")
	// Set up a channel to listen for OS signals to gracefully exit
	// This allows us to clean up resources and unsubscribe before exiting
	sigChan := make(chan os.Signal, 1) // Create a channel to listen for OS signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM) // Listen for interrupt signals (Ctrl+C)
	<-sigChan // Wait for a signal to exit

	// Unsubscribe and Disconnect
	client.Unsubscribe(topic)
	client.Disconnect(250)
	fmt.Println("Subscriber disconnected.")
}