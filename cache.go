package omgo

import (
	"sync"
	"time"
)

type cacheItem struct {
	data       []byte
	expiration time.Time
}

type Cache struct {
	items map[string]cacheItem
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]cacheItem),
	}
}

func (c *Cache) Set(key string, data []byte, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		data:       data,
		expiration: time.Now().Add(expiration),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	if time.Now().After(item.expiration) {
		delete(c.items, key)
		return nil, false
	}
	return item.data, true
}
