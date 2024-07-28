package pkg

import (
	"sync"
	"time"
)

//CacheItem 表示緩存中的一個項目
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}
	
// Cache 表示緩存結構
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}
	
// NewCache 創建一個新的 Cache
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}
	
// Get 獲取緩存中的項目
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || item.ExpiresAt.Before(time.Now()) {
		return nil, false
	}
	return item.Data, true
}

// Set 設置緩存中的項目
func (c *Cache) Set(key string, data interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(duration),
	}
}