
# --- Stage 1: The Builder ---
FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../ ./

# build the 'location-service' command and name the output binary '/location-service'.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /location-service ./cmd/location-service

# --- Stage 2: The Final Image ---
FROM alpine:latest

WORKDIR /

COPY --from=builder /location-service /location-service

CMD ["/location-service"]