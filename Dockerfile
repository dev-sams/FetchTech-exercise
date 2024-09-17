# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create a lightweight image to run the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Expose the port that Gin will run on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
