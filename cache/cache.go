package cache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]*CacheItem
}

type CacheItem struct {
	value []byte
	ttl   time.Duration
}

func New() *Cache {
	return &Cache{
		mu:   sync.RWMutex{},
		data: make(map[string]*CacheItem, 0),
	}
}

func (c *Cache) Set(key []byte, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[string(key)] = &CacheItem{
		value: value,
		ttl:   ttl,
	}
	return nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[string(key)]
	if !ok {
		return nil, errors.New("item not found")
	}
	return value.value, nil
}

func (c *Cache) Delete(key []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, string(key))
	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.data[string(key)]
	return ok
}
