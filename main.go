package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type URLShortener struct {
	urls          map[string]string
	domainMetrics map[string]int
	mu            sync.Mutex
}

func (us *URLShortener) shortenURL(w http.ResponseWriter, r *http.Request) {
	var input struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.URL == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	if shortURL, exists := us.urls[input.URL]; exists {
		fmt.Fprintf(w, `{"shortURL": "%s"}`, shortURL)
		return
	}

	hash := generateShortURL(input.URL)
	us.urls[input.URL] = hash
	us.domainMetrics[getDomain(input.URL)]++

	fmt.Fprintf(w, `{"shortURL": "%s"}`, hash)
}

func (us *URLShortener) redirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	us.mu.Lock()
	defer us.mu.Unlock()

	if originalURL, exists := us.getOriginalURL(shortURL); exists {
		http.Redirect(w, r, originalURL, http.StatusFound)
		return
	}

	http.Error(w, "URL not found", http.StatusNotFound)
}

func (us *URLShortener) getTopDomains(w http.ResponseWriter, r *http.Request) {
	us.mu.Lock()
	defer us.mu.Unlock()

	topDomains := us.getTopDomainsList(3)
	json.NewEncoder(w).Encode(topDomains)
}

func generateShortURL(url string) string {
	hasher := sha1.New()
	hasher.Write([]byte(url))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:6]
}

func getDomain(urlString string) string {
	parsedURL, _ := url.Parse(urlString)
	return parsedURL.Hostname()
}

func (us *URLShortener) getOriginalURL(shortURL string) (string, bool) {
	for originalURL, hashedURL := range us.urls {
		if hashedURL == shortURL {
			return originalURL, true
		}
	}
	return "", false
}

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

	http.HandleFunc("/shorten", shortener.shortenURL)
	http.HandleFunc("/", shortener.redirectToOriginalURL)
	http.HandleFunc("/metrics/top-domains", shortener.getTopDomains)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
