package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value       any
	ExpiredTime time.Time
}

type Cache struct {
	mu   sync.Mutex
	data map[string]CacheItem
}

func NewCache(CleanUpInterval time.Duration) *Cache {
	cache := &Cache{data: make(map[string]CacheItem)}
	cache.goCleanup(CleanUpInterval)

	return cache
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheItem{Value: value, ExpiredTime: time.Now().Add(ttl)}
}

func (c *Cache) Get(key string) any {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.data[key]
	if !ok {
		return nil
	}

	if time.Now().After(item.ExpiredTime) {
		delete(c.data, key)
		return nil
	}

	return item.Value
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
