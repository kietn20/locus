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
	DBConn *pgx.Conn
	MQTTClient mqtt.Client
	vehicleStates map[string]string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// --- Connect to Database ---
	conn := db.Connect()
	defer conn.Close(context.Background())

	// --- Connect to MQTT ---
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("geofence-service")
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}
	fmt.Println("Geofence service connected to MQTT broker.")

	service := &GeofenceService{
		DBConn: conn,
		MQTTClient: mqttClient,
		vehicleStates: make(map[string]string),
	}

	opts.SetDefaultPublishHandler(service.messageHandler)

	topic := "locus/hehicles/+/location"
	if token := mqttClient.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe:  %v", token.Error())
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

func (s *GeofenceService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	var loc vehicle.LocationData
	if err := loc.FromJSON(msg.Payload()); err != nil {
		log.Printf("Error unmarshalling location data: %v", err)
		return
	}



	fmt.Printf("Processing location for %s\n", loc.VehicleID)
}