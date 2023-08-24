# Go URL Shortener

This is a simple URL shortener application built using the Go programming language.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Dockerization](#dockerization)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites

- Go (1.17 or higher)
- Docker (optional, for containerization)

## Installation

1. Clone the repository to your local machine:

    ```bash
    git clone https://github.com/gordiancyber/url-shortener-app.git
    cd url-shortener-app
    ```

2. Install the project dependencies:

    ```bash
    go mod download
    ```

## Usage

1. Run the application:

    ```bash
    go run main.go
    ```

2. To shorten a URL, make a POST request:

    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"url": "https://www.example.com"}' http://localhost:8080/shorten
    ```

3. To access the original URL using the short URL:

    ```bash
    curl http://localhost:8080/<short_url>
    ```
4. To retrieve the top domain metrics:
   
   ```bash
   curl http://localhost:8080/metrics/top-domains
   ```
   
## Dockerization

1. Build the Docker image:

    ```bash
    docker build -t url-shortener-app .
    ```

2. Run the Docker container:

    ```bash
    docker run -p 8080:8080 url-shortener-app
    ```
    
3. Pull Docker Image here :
   ```bash
   docker pull heygordian/url-shortener-app:latest
   ```
   ## API Endpoints
- `POST /shorten`: Shorten a URL.
  Example: 
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"url": "https://www.example.com"}' http://localhost:8080/shorten

