# Project Locus: Real-Time Fleet Tracking & Geofencing System

## High-Level Summary

Project Locus is a complete, event-driven backend system designed to track a fleet of vehicles in real-time and generate alerts when they enter or exit user-defined geographic boundaries (geofences). The system is built using a modern microservices architecture, is fully containerized with Docker, and is deployed on AWS. This project demonstrates a deep understanding of Go, concurrency, real-time messaging with MQTT, REST API design, and cloud infrastructure management.

---
## Architecture Diagram

The system consists of several independent Go services that communicate asynchronously via an MQTT message broker and a shared PostgreSQL database.






---

## Tech Stack

*   **Language:** Go
*   **Messaging:** MQTT (Eclipse Mosquitto Broker)
*   **Database:** PostgreSQL with the PostGIS extension for geospatial queries.
*   **API:** RESTful API built with the `chi` router.
*   **Infrastructure:** Docker, Docker Compose
*   **Cloud Provider:** Amazon Web Services (AWS)
    *   **Compute:** EC2 (`t2.micro`)
    *   **Networking:** VPC, Security Groups

---