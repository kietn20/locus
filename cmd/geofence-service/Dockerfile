
FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY ../../go.mod ../../go.sum ./

RUN go mod download

COPY ../../ ./


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /geofence-service ./cmd/geofence-service

FROM alpine:latest

WORKDIR /

COPY --from=builder /geofence-service /geofence-service

CMD ["/geofence-service"]
