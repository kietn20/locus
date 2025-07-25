# Project Locus: Real-Time Fleet Tracking & Geofencing System

## High-Level Summary

Project Locus is a complete, event-driven backend system designed to track a fleet of vehicles in real-time and generate alerts when they enter or exit user-defined geographic boundaries (geofences). The system is built using a modern microservices architecture, is fully containerized with Docker, and is deployed on AWS. This project demonstrates a deep understanding of Go, concurrency, real-time messaging with MQTT, REST API design, and cloud infrastructure management.

## Architecture Diagram

The system consists of several independent Go services that communicate asynchronously via an MQTT message broker and a shared PostgreSQL database.

![Architecture Diagram](https://github.com/kietn20/locus/blob/main/diagram.png)

## Demo

![Demo Image](https://github.com/kietn20/locus/blob/main/demo-img.png)

Youtube Demo: [LINK](https://www.youtube.com/watch?v=ldhcy7bkz5E)

## Tech Stack

-   **Language:** Go
-   **Messaging:** MQTT (Eclipse Mosquitto Broker)
-   **Database:** PostgreSQL with the PostGIS extension for geospatial queries.
-   **API:** RESTful API built with the `chi` router.
-   **Infrastructure:** Docker, Docker Compose
-   **Cloud Provider:** Amazon Web Services (AWS)
    -   **Compute:** EC2 (`t2.micro`)
    -   **Networking:** VPC, Security Groups

---

## Features

-   **Real-Time Location Ingestion:** The `location-service` subscribes to MQTT topics and persists vehicle location data to the database.
-   **Dynamic Geofence Management:** A REST API (`api-service`) allows for creating, viewing, and deleting polygonal geofences.
-   **Stateful Event Engine:** The `geofence-service` tracks the state of each vehicle, using PostGIS to perform efficient geospatial calculations and detect when a vehicle enters or exits a geofence.
-   **Event-Driven Alerts:** Generates new MQTT messages on a separate topic (`locus/geofence/events`) for every detected enter/exit event.
-   **Scalable Simulation:** The `vehicle-simulator` uses goroutines to simulate a configurable number of vehicles concurrently.
-   **Cloud Deployed:** The entire backend stack is containerized and deployed to an AWS EC2 instance.

---

## How to Run

### Local Development

1.  **Prerequisites:** Go, Docker, Docker Compose installed.
2.  Clone the repository: `git clone ...`
3.  Create a `.env` file from the project root (see `.env.example` if available).
4.  Start the infrastructure: `docker-compose up -d`
5.  Run the services in separate terminals:
    ```bash
    go run ./cmd/api-service/main.go
    go run ./cmd/location-service/main.go
    go run ./cmd/geofence-service/main.go
    ```
6.  Create a geofence (see API Usage section).
7.  Run the simulator: `go run ./cmd/vehicle-simulator/main.go`

### Cloud Deployment (AWS)

The application is deployed on an EC2 instance and managed via Docker Compose.

1.  Set up an EC2 instance with Docker, Docker Compose, and Git installed.
2.  Configure the security group to allow inbound traffic on ports 22 (SSH), 80 (HTTP), and 1883 (MQTT).
3.  Clone the repository onto the instance.
4.  Create the `.env` file on the server.
5.  Run the backend stack: `sudo docker-compose -f docker-compose.deploy.yml up -d`
6.  To test, run the vehicle simulator from a local machine with the `MQTT_BROKER_HOST` environment variable set to the EC2 instance's public IP.

---

## API Usage

### Create a Geofence

-   **Endpoint:** `POST /api/v1/geofences`
-   **Method:** `POST`
-   **Body (raw JSON):**
    ```json
    {
    	"type": "Feature",
    	"properties": { "name": "downtown-la" },
    	"geometry": {
    		"type": "Polygon",
    		"coordinates": [
    			[
    				[-118.25, 34.04],
    				[-118.23, 34.04],
    				[-118.23, 34.06],
    				[-118.25, 34.06],
    				[-118.25, 34.04]
    			]
    		]
    	}
    }
    ```

### List Geofences

-   **Endpoint:** `GET /api/v1/geofences`
-   **Method:** `GET`
-   **Returns:** A JSON array of all geofences in the database.
