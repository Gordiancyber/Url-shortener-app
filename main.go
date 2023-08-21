package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	rdb           *redis.Client
	originalToShortURL = make(map[string]string)
	domainToCount = make(map[string]int)
	lock sync.Mutex
	ctx = context.Background()
)

// NewRedisClient creates and returns a new Redis client instance.
func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}

func main() {
	fmt.Println("URL Shortener with Go and Redis")

	// Initialize the Redis client
	rdb = NewRedisClient()
	defer rdb.Close()

	// Handle the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the URL Shortener!")
	})

	// Handle the /shorten endpoint
	http.HandleFunc("/shorten", shortenHandler)

	// Handle the /r/ endpoint for redirection
	http.HandleFunc("/r/", redirectHandler)

	// Handle the /metrics endpoint
	http.HandleFunc("/metrics", metricsHandler)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// shortenHandler handles shortening URLs and preventing duplicates.
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")

		// Ensure thread-safe access to data
		lock.Lock()
		defer lock.Unlock()

		// Check if the original URL has already been shortened
		if shortURL, ok := originalToShortURL[url]; ok {
			jsonResponse, _ := json.Marshal(map[string]string{"short_url": shortURL})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponse)
			return
		}

		// Generate a new short URL and store the mapping
		shortURL, err := shortenURL(rdb, url)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Store the mapping of original URL to short URL
		originalToShortURL[url] = shortURL

		// Extract the domain from the input URL
		u, err := url.Parse(url)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		domain := strings.TrimPrefix(u.Host, "www.")
		domainToCount[domain]++

		// Respond with the short URL
		jsonResponse, _ := json.Marshal(map[string]string{"short_url": shortURL})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}

// redirectHandler handles redirecting short URLs to original URLs.
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/r/"):]
	url, err := redirectToURL(rdb, shortURL)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

// metricsHandler handles providing the top 3 domains with the highest counts.
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure thread-safe access to data
	lock.Lock()
	defer lock.Unlock()

	// Create a slice of domain-count pairs for sorting
	var domainCountPairs []struct {
		Domain string
		Count  int
	}
	for domain, count := range domainToCount {
		domainCountPairs = append(domainCountPairs, struct {
			Domain string
			Count  int
		}{Domain: domain, Count: count})
	}

	// Sort the domain-count pairs by count in descending order
	sort.Slice(domainCountPairs, func(i, j int) bool {
		return domainCountPairs[i].Count > domainCountPairs[j].Count
	})

	// Get the top 3 domains
	topDomains := domainCountPairs[:3]

	// Prepare the JSON response
	response := make(map[string]int)
	for _, pair := range topDomains {
		response[pair.Domain] = pair.Count
	}

	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
