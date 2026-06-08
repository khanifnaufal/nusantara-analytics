package cache

import (
	"os"
	"strconv"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var c *gocache.Cache

// Init initializes the global cache singleton instance
func Init() {
	defaultTTLSec := 300
	if ttlStr := os.Getenv("CACHE_DEFAULT_TTL"); ttlStr != "" {
		if val, err := strconv.Atoi(ttlStr); err == nil {
			defaultTTLSec = val
		}
	}

	defaultTTL := time.Duration(defaultTTLSec) * time.Second
	cleanupInterval := 10 * time.Minute

	c = gocache.New(defaultTTL, cleanupInterval)
}

// Get retrieves an item from the cache.
// Returns the value and a boolean indicating whether the key was found.
func Get(key string) (interface{}, bool) {
	return c.Get(key)
}

// Set adds an item to the cache, replacing any existing item.
// If the duration is 0 (gocache.DefaultExpiration / time.Duration(0)), the cache's default expiration is used.
// If the duration is -1 (gocache.NoExpiration / time.Duration(-1)), the item never expires.
func Set(key string, value interface{}, ttl time.Duration) {
	c.Set(key, value, ttl)
}

// Delete removes an item from the cache.
func Delete(key string) {
	c.Delete(key)
}
