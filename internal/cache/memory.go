package cache

import "sync"

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewMemory() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]Item),
	}
}

func (c *MemoryCache) Get(key string) (Item, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	return item, exists
}

func (c *MemoryCache) Set(key string, item Item) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item
}

func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]Item)
}
