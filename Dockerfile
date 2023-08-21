# Use the official Go image as the base image
FROM golang:1.17-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module and Go sum files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go application
RUN go build -o url-shortener-app .

# Expose the port your application listens on
EXPOSE 8080

# Start the application when the container starts
CMD ["./url-shortener-app"]
