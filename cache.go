package simplecache

import (
	"time"
)

type Item struct {
	expiry time.Duration
	value  interface{}
	key    string
}

type cacher interface {
	len() int
	get(key string) (interface{}, bool)
	set(key string, val interface{}, expiry time.Duration) bool
	remove(key string) bool
	purge() bool
}

type Cache struct {
	cache cacher
}

func New(maxItems int64) *Cache {
	cache := newLRU(100)
	return &Cache{cache: cache}
}

func (c *Cache) Set(key string, value interface{}, expiry time.Duration) bool {
	return c.cache.set(key, value, expiry)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.cache.get(key)
}

func (c *Cache) Purge() {
	c.cache.purge()
}

func (c *Cache) Len() int {
	return c.cache.len()
}
