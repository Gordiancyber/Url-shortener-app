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



