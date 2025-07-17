# --- Stage 1: The Builder ---
# start with the official Go image, which contains the Go compiler.
# pinning it to a specific version for reproducibility.
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the Go module files. This is done as a separate layer to leverage Docker's
# layer caching. If go.mod and go.sum don't change, Docker won't re-download dependencies.
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy the rest of the source code.
# The ../../ is because this Dockerfile is two levels deep from the project root.
COPY ../../ ./


# Compile the application.
# The -o flag specifies the output file name.
# build for linux/amd64 since that's what most Docker hosts (and cloud servers) run.
# CGO_ENABLED=0 creates a statically linked binary (no external C libraries needed).
# -ldflags="-w -s" strips debug information, making the binary smaller.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /api-service ./cmd/api-service

# --- Stage 2: The Final Image ---
# start from a "scratch" image, which is completely empty. It's the most minimal possible.
FROM scratch

# Set the working directory.
WORKDIR /

# Copy ONLY the compiled binary from the 'builder' stage.
COPY --from=builder /api-service /api-service

# Set the command to run when the container starts.
# This tells the container to execute our compiled Go program.
CMD ["/api-service"]