package main // Package declaration for the main program

import (
	"crypto/sha1"         // Importing the crypto/sha1 package for SHA-1 hashing
	"encoding/hex"        // Importing the encoding/hex package for hexadecimal encoding
	"encoding/json"       // Importing the encoding/json package for JSON handling
	"fmt"                 // Importing the fmt package for formatted I/O
	"net/http"            // Importing the net/http package for HTTP server functionality
	"net/url"             // Importing the net/url package for URL parsing
	"strings"             // Importing the strings package for string manipulation
	"sync"                // Importing the sync package for synchronization utilities
)

// URLShortener is a struct representing the URL shortening service.
type URLShortener struct {
	urls          map[string]string // Mapping from original URLs to short URLs
	domainMetrics map[string]int    // Metrics for domains (e.g., how many times a domain is shortened)
	mu            sync.Mutex        // Mutex for concurrent access to data structures
}

// shortenURL handles the shortening of a given URL.
func (us *URLShortener) shortenURL(w http.ResponseWriter, r *http.Request) {
	var input struct {
		URL string `json:"url"` // Struct for decoding JSON request body
	}

	// Decode JSON request and check for errors or empty input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.URL == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest) // Return a bad request error
		return
	}

	// Lock to ensure thread safety when accessing data structures
	us.mu.Lock()
	defer us.mu.Unlock()

	// Check if the URL is already shortened, if so, return the existing short URL
	if shortURL, exists := us.urls[input.URL]; exists {
		fmt.Fprintf(w, `{"shortURL": "%s"}`, shortURL) // Respond with the existing short URL
		return
	}

	// Generate a short URL hash and store it in the mappings
	hash := generateShortURL(input.URL)
	us.urls[input.URL] = hash
	us.domainMetrics[getDomain(input.URL)]++

	// Respond with the generated short URL
	fmt.Fprintf(w, `{"shortURL": "%s"}`, hash)
}

// redirectToOriginalURL redirects to the original URL corresponding to a short URL.
func (us *URLShortener) redirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")

	us.mu.Lock()
	defer us.mu.Unlock()

	// Look up the original URL and redirect if found, otherwise, return an error
	if originalURL, exists := us.getOriginalURL(shortURL); exists {
		http.Redirect(w, r, originalURL, http.StatusFound) // Redirect to the original URL
		return
	}

	http.Error(w, "URL not found", http.StatusNotFound) // Return a not found error
}

// getTopDomains returns a JSON list of the top domains along with their metrics.
func (us *URLShortener) getTopDomains(w http.ResponseWriter, r *http.Request) {
	us.mu.Lock()
	defer us.mu.Unlock()

	// Get the list of top domains and encode them as JSON
	topDomains := us.getTopDomainsList(3)
	json.NewEncoder(w).Encode(topDomains)
}

// generateShortURL generates a short URL hash using SHA-1 and returns the first 6 characters.
func generateShortURL(url string) string {
	hasher := sha1.New()
	hasher.Write([]byte(url))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:6]
}

// getDomain extracts the domain from a given URL.
func getDomain(urlString string) string {
	parsedURL, _ := url.Parse(urlString)
	return parsedURL.Hostname()
}

// getOriginalURL returns the original URL given a short URL.
func (us *URLShortener) getOriginalURL(shortURL string) (string, bool) {
	for originalURL, hashedURL := range us.urls {
		if hashedURL == shortURL {
			return originalURL, true
		}
	}
	return "", false
}

// getTopDomainsList returns a map of the top domains and their metrics, limited to the specified limit.
func (us *URLShortener) getTopDomainsList(limit int) map[string]int {
	domains := make([]string, 0, len(us.domainMetrics))
	for domain := range us.domainMetrics {
		domains = append(domains, domain)
	}

	topDomains := make(map[string]int, limit)
	for _, domain := range domains {
		topDomains[domain] = us.domainMetrics[domain]
	}

	return topDomains
}

func main() {
	shortener := &URLShortener{
		urls:          make(map[string]string),
		domainMetrics: make(map[string]int),
	}

	// Define HTTP routes and handlers
	http.HandleFunc("/shorten", shortener.shortenURL)
	http.HandleFunc("/", shortener.redirectToOriginalURL)
	http.HandleFunc("/metrics/top-domains", shortener.getTopDomains)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil) // Start the HTTP server on port 8080
}
