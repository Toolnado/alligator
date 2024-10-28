package cache

import (
	"errors"
	"log"
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
	log.Printf("SET command: key=%s, value=%s, ttl=%s\n", key, value, ttl)
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
	log.Printf("GET command: key=%s\n", key)
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, errors.New("item not found")
	}
	return value.value, nil
}

func (c *Cache) Delete(key string) error {
	log.Printf("DELETE command: key=%s\n", key)
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *Cache) Has(key string) bool {
	log.Printf("HAS command: key=%s\n", key)
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.data[key]
	return ok
}
