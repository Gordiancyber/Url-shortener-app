package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8" // Import the Redis client library
	"golang.org/x/net/context"
)

var (
	ctx        = context.Background()
	redisClient *redis.Client

	mutex      sync.Mutex
)

func generateShortURL(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))[:6]
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	// ...
}

func getShortURL(originalURL string) (string, bool) {
	// ...
}

func storeURLMapping(shortURL, originalURL string) {
	// ...
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// ...
}

func main() {
	// Initialize the Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",              // No password
		DB:       0,               // Default DB
	})

	// Handle URL shortening requests
	http.HandleFunc("/shorten", shortenURLHandler)

	// Handle short URL redirection
	http.HandleFunc("/", redirectHandler)

	// Start the server on port 80
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}


