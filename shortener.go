package main

import (
    "crypto/sha1"
    "encoding/base64"
    "github.com/go-redis/redis/v8"
)

// Define a constant prefix for the Redis keys
const keyPrefix = "urlshortener:"

// shortenURL generates a short URL and stores the mapping in Redis.
func shortenURL(rdb *redis.Client, url string) (string, error) {
    // Calculate the SHA-1 hash of the URL
    hash := sha1.New()
    hash.Write([]byte(url))
    // Encode the hash using base64 to get a short URL
    shortURL := base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]

    // Store the mapping of short URL to original URL in Redis
    err := rdb.Set(ctx, keyPrefix+shortURL, url, 0).Err()
    if err != nil {
        return "", err
    }

    // Return the generated short URL
    return shortURL, nil
}

// redirectToURL retrieves the original URL using the short URL.
func redirectToURL(rdb *redis.Client, shortURL string) (string, error) {
    // Retrieve the original URL from Redis using the short URL as the key
    url, err := rdb.Get(ctx, keyPrefix+shortURL).Result()
    if err != nil {
        return "", err
    }

    // Return the original URL
    return url, nil
}
