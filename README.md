# Project Locus: Real-Time Fleet Tracking & Geofencing System

## High-Level Summary

Project Locus is a complete, event-driven backend system designed to track a fleet of vehicles in real-time and generate alerts when they enter or exit user-defined geographic boundaries (geofences). The system is built using a modern microservices architecture, is fully containerized with Docker, and is deployed on AWS. This project demonstrates a deep understanding of Go, concurrency, real-time messaging with MQTT, REST API design, and cloud infrastructure management.

---
## Architecture Diagram

The system consists of several independent Go services that communicate asynchronously via an MQTT message broker and a shared PostgreSQL database.