package main

import (
    "context"
    "github.com/go-redis/redis/v8"
)

// Create a context.Background() instance to be used globally
var ctx = context.Background()

// NewRedisClient initializes and returns a new Redis client instance.
func NewRedisClient() *redis.Client {
    // Create a new Redis client with the specified options
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Address of the Redis server
        Password: "",               // Redis server password (empty if not required)
        DB:       0,                // Database index to use
    })

    // Return the
