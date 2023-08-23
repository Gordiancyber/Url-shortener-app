FROM golang:1.21

# Set the working directory within the container
WORKDIR /url-shortener-app

# Copy the entire project to the container's working directory
COPY . /url-shortener-app

# Build the Go application
RUN go build -o url-shortener-app .



