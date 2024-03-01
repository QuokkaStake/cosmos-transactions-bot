package cache

import (
	"time"

	"github.com/rs/zerolog"
)

type CacheEntry struct {
	Value    interface{}
	StoredAt time.Time
}

type Cache struct {
	Logger    zerolog.Logger
	StoreTime time.Duration
	Entries   map[string]CacheEntry
}

func NewCache() *Cache {
	return &Cache{
		StoreTime: 10 * time.Minute,
		Entries:   make(map[string]CacheEntry),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	entry, found := c.Entries[key]
	if !found {
		return nil, false
	}

	// cache expired
	if entry.StoredAt.Add(c.StoreTime).Before(time.Now()) {
		return nil, false
	}

	return entry.Value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.Entries[key] = CacheEntry{
		Value:    value,
		StoredAt: time.Now(),
	}
}
