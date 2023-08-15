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




