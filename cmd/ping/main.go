package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Define MQTT broker connection options
	// "tcp://localhost:1883" points to the Mosquitto container we are running
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("go-publisher") // Each client must have a unique ID

	// Create a new client
	client := mqtt.NewClient(opts)

	// Connect to the broker
	// The .Token variants block until the action is complete and return an error.
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %v", token.Error())
	}
	fmt.Println("Publisher connected to broker.")

	// The topic we are publishing to
	topic := "locus/test"
	message := "Hello from Go Publisher!"

	// Publish the message
	token := client.Publish(topic, 0, false, message)
	token.Wait() // Wait for the publish to complete

	fmt.Printf("Published message '%s' to topic '%s'\n", message, topic)

	// Disconnect from the broker
	client.Disconnect(250)
	fmt.Println("Publisher disconnected.")
}