# Use a minimal base image with only necessary dependencies
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy only the necessary files for dependency resolution
COPY go.mod go.sum ./

# Download and cache GO dependencies
RUN go mod download

# Copy the entire application source code
COPY . /app

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Create a minimal production image
FROM alpine:latest

WORKDIR /app

# Copy only the compile binary from the builder image
COPY --from=builder /app/app .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./app"]

