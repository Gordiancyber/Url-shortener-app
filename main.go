package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"sync"
)

type URLShortener struct {
	storage     map[string]string
	accessCount map[string]int
	mu          sync.Mutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		storage:     make(map[string]string),
		accessCount: make(map[string]int),
		mu:          sync.Mutex{},
	}
}

func (s *URLShortener) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	shortURL, found := s.storage[input.URL]
	if !found {
		shortURL = generateShortURL(input.URL)
		s.storage[input.URL] = shortURL
	}

	response := struct {
		ShortURL string `json:"short_url"`
	}{
		ShortURL: shortURL,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *URLShortener) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]

	s.mu.Lock()
	defer s.mu.Unlock()

	originalURL, found := s.storage[shortURL]
	if !found {
		http.NotFound(w, r)
		return
	}

	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
		return
	}

	domain := parsedURL.Host
	s.accessCount[originalURL]++
	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}

func (s *URLShortener) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	type DomainMetric struct {
		Domain string `json:"domain"`
		Count  int    `json:"count"`
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	domainCount := make(map[string]int)
	for url := range s.accessCount {
		parsedURL, err := url.Parse(url)
		if err != nil {
			continue
		}
		domain := parsedURL.Host
		domainCount[domain]++
	}

	var metrics []DomainMetric
	for domain, count := range domainCount {
		metrics = append(metrics, DomainMetric{Domain: domain, Count: count})
	}

	sortMetrics(metrics)

	if len(metrics) > 3 {
		metrics = metrics[:3]
	}

	json.NewEncoder(w).Encode(metrics)
}

func generateShortURL(url string) string {
	// In a real system, you would generate a unique short URL using an algorithm
	// For simplicity, let's just use a hash of the URL here
	return fmt.Sprintf("short_%d", len(url))
}

func sortMetrics(metrics []DomainMetric) {
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Count > metrics[j].Count
	})
}

func main() {
	shortener := NewURLShortener()

	http.HandleFunc("/shorten", shortener.ShortenHandler)
	http.HandleFunc("/", shortener.RedirectHandler)
	http.HandleFunc("/metrics", shortener.MetricsHandler)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


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


