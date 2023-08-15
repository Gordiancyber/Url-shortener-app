package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8" // Import the Redis client library
	"golang.org/x/net/context"
)

