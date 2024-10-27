package cache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]*Item
}

type Item struct {
	value []byte
	ttl   time.Duration
}

func New() *Cache {
	return &Cache{
		mu:   sync.RWMutex{},
		data: make(map[string]*Item),
	}
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = &Item{
		value: value,
		ttl:   ttl,
	}
	go c.deleteByTTL(ttl, key)
	return nil
}

func (c *Cache) deleteByTTL(ttl time.Duration, key string) {
	<-time.After(ttl)
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, errors.New("item not found")
	}
	return value.value, nil
}

func (c *Cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *Cache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.data[key]
	return ok
}
