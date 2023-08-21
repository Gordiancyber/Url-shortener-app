package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

func main() {
    fmt.Println("URL Shortener with Go and Redis")

    rdb := NewRedisClient()
    defer rdb.Close()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to the URL Shortener!")
    })

    http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            url := r.FormValue("url")
            shortURL, err := shortenURL(rdb, url)
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            jsonResponse, _ := json.Marshal(map[string]string{"short_url": shortURL})
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusCreated)
            w.Write(jsonResponse)
        }
    })

    http.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
        shortURL := r.URL.Path[len("/r/"):]
        url, err := redirectToURL(rdb, shortURL)
        if err != nil {
            http.Error(w, "Not Found", http.StatusNotFound)
            return
        }

        http.Redirect(w, r, url, http.StatusSeeOther)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}