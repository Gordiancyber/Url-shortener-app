package main

import (
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

    log.Fatal(http.ListenAndServe(":8080", nil))
}
