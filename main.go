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

	rdb = NewRedisClient()
	defer rdb.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the URL Shortener!")
	})

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/r/", redirectHandler)
	http.HandleFunc("/metrics", metricsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")

		lock.Lock()
		defer lock.Unlock()

		if shortURL, ok := originalToShortURL[url]; ok {
			jsonResponse, _ := json.Marshal(map[string]string{"short_url": shortURL})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponse)
			return
		}

		shortURL, err := shortenURL(rdb, url)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		originalToShortURL[url] = shortURL

		u, err := url.Parse(url)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		domain := strings.TrimPrefix(u.Host, "www.")
		domainToCount[domain]++

		jsonResponse, _ := json.Marshal(map[string]string{"short_url": shortURL})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/r/"):]
	url, err := redirectToURL(rdb, shortURL)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

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

	sort.Slice(domainCountPairs, func(i, j int) bool {
		return domainCountPairs[i].Count > domainCountPairs[j].Count
	})

	topDomains := domainCountPairs[:3]

	response := make(map[string]int)
	for _, pair := range topDomains {
		response[pair.Domain] = pair.Count
	}

	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
