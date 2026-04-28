package cache

import (
	"arc-events/models"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	items models.ArcDataResponse
}

var (
	instance *Cache
	once     sync.Once
)

func (c *Cache) Set(value models.ArcDataResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = value
}

func (c *Cache) Get() models.ArcDataResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.items
}

func GetInstance() *Cache {
	once.Do(func() {
		instance = &Cache{
			items: models.ArcDataResponse{},
		}
	})
	return instance
}
