# Stage 1: Build the Go application
# Using a version of Go that is compatible with your go.mod file
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go module files to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application, creating a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /crm-backend .

# Stage 2: Create a small, secure final image
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the built binary from the 'builder' stage
COPY --from=builder /crm-backend .

# Expose port 8080 to the outside world. This is the port your Go app listens on.
EXPOSE 8080

# Command to run the executable when the container starts
CMD ["./crm-backend"]