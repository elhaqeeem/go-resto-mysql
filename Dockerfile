# Stage 1: Build the Go application
FROM golang:1.21-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Copy CA certificate
COPY certs/ca.pem /etc/ssl/certs/ca.pem

# Build the Go application as a static binary
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o main .

# Stage 2: Create a minimal image for the application
FROM scratch

WORKDIR /root/

# Copy the compiled static binary from the builder stage
COPY --from=builder /app/main .

EXPOSE 8080
# Command to run the executable
CMD ["./main"]
