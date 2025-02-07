FROM golang:1.23 AS builder
WORKDIR /backend

COPY go.mod go.sum /backend/
RUN go mod tidy

# Copy the entire project
COPY . /backend/

# Enable Go modules
ENV GO111MODULE=on

# Debug: Verify main.go exists
RUN ls -la /backend/cmd/server  # Debug step

# Build the application inside the correct directory
WORKDIR /backend/cmd/server
RUN go build -o /backend/main .  # <-- Build to /backend/main

# Final container using Golang image
FROM golang:1.23
WORKDIR /root/

# Copy the compiled binary from the builder
COPY --from=builder /backend/main .

# Copy the entire backend project
COPY --from=builder /backend /root/backend

# Copy .env file
COPY .env .env

# Set correct environment variables
ENV GO111MODULE=on
WORKDIR /root/backend

CMD ["./main"]
