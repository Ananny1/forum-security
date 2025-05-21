# Use the official Golang image
FROM golang:1.21.6

# Set the Current Working Directory inside the container
WORKDIR /app

# Add metadata labels
LABEL authors="Mariam, Fatima, Reem, Sayed Hassan, Mohamed"
LABEL description="A Docker image for building and running a Go application"

# Copy the Go Modules files
COPY go.mod .

# Download dependencies as separate layers to improve caching
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o bin/server .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/app/bin/server"]

#openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt
