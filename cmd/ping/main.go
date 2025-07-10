package main

import (
	"fmt"
	"log"
	// "time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("go-publisher")

	client := mqtt.NewClient(opts)


	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %v", token.Error())
	}
	fmt.Println("Publisher connected to broker.")

	topic := "locus/test"
	message := "Hello from Go Publisher!"

	token := client.Publish(topic, 0, false, message)
	token.Wait()

	fmt.Printf("Published message '%s' to topic '%s'\n", message, topic)

	client.Disconnect(250)
	fmt.Println("Publisher disconnected.")
}