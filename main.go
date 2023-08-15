package main

import (
	"crypto/md5"          // Import the md5 package to generate hashes
	"encoding/hex"        // Import the hex package to encode hashes as strings
	"encoding/json"       // Import the json package for JSON operations
	"fmt"                // Import the fmt package for printing
	"net/http"           // Import the http package for building HTTP servers
	"sync"               // Import the sync package for handling concurrency
)

var (
	urlMapping = make(map[string]string) // A map to store short URLs and their corresponding original URLs
	mutex      sync.Mutex                // A mutex to ensure thread-safe access to the urlMapping
)

func generateShortURL(url string) string {
	hasher := md5.New()                      // Create an MD5 hash object
	hasher.Write([]byte(url))                // Write the original URL bytes to the hash object
	return hex.EncodeToString(hasher.Sum(nil))[:6] // Generate a hex-encoded hash and truncate to 6 characters
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"` // Struct to hold the incoming JSON data
	}
	err := json.NewDecoder(r.Body).Decode(&data) // Decode JSON request body into the data struct
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if data.URL == "" {
		http.Error(w, "URL missing", http.StatusBadRequest)
		return
	}

	mutex.Lock()         // Lock the mutex to ensure thread-safe access to urlMapping
	defer mutex.Unlock() // Unlock the mutex when the function exits

	shortURL := generateShortURL(data.URL) // Generate a short URL for the original URL
	urlMapping[shortURL] = data.URL        // Store the original URL in the urlMapping

	response := map[string]string{"short_url": shortURL}         // Create a response map
	responseJSON, _ := json.Marshal(response)                    // Convert the response map to JSON
	w.Header().Set("Content-Type", "application/json")            // Set the response header
	w.WriteHeader(http.StatusCreated)                             // Set the HTTP status code
	w.Write(responseJSON)                                        // Write the JSON response
}


