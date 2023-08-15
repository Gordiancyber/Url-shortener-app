package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	urlMapping = make(map[string]string)
	mutex      sync.Mutex
)

func generateShortURL(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil))[:6]
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if data.URL == "" {
		http.Error(w, "URL missing", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	shortURL := generateShortURL(data.URL)
	urlMapping[shortURL] = data.URL

	response := map[string]string{"short_url": shortURL}
	responseJSON, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	if originalURL, exists := urlMapping[shortURL]; exists {
		http.Redirect(w, r, originalURL, http.StatusFound)
	} else {
		http.Error(w, "Short URL not found", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
