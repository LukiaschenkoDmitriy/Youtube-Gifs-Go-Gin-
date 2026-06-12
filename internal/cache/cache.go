package cache

import (
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	Value       any
	ExpiredTime time.Time
}

type Cache struct {
	mu    sync.Mutex
	data  map[string]CacheItem
	debug bool
}

func NewCache(CleanUpInterval time.Duration, debug bool) *Cache {
	cache := &Cache{data: make(map[string]CacheItem), debug: debug}
	cache.goCleanup(CleanUpInterval)

	return cache
}

func (c *Cache) debugWrite(key string, value any, ttl time.Duration) {
	fmt.Printf("\033[32m[Cache Write]\033[0m\n  Key   : %s\n  Value : %v\n  TTL   : %s\n",
		key, value, ttl)
}

func (c *Cache) debugRead(key string, value any, ttl time.Duration) {
	fmt.Printf("\033[36m[Cache Read]\033[0m\n  Key   : %s\n  Value : %v\n  TTL   : %s\n",
		key, value, ttl)
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheItem{Value: value, ExpiredTime: time.Now().Add(ttl)}
	if c.debug {
		c.debugWrite(key, value, ttl)
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(item.ExpiredTime) {
		delete(c.data, key)
		return nil, false
	}

	if c.debug {
		c.debugRead(key, item.Value, time.Until(item.ExpiredTime))
	}

	return item.Value, true
}

func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Call function once
func (c *Cache) goCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			c.mu.Lock()
			for key, item := range c.data {
				if time.Now().After(item.ExpiredTime) {
					delete(c.data, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}
