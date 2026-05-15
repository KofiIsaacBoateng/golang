package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]CacheMap
	mu *sync.Mutex
}

type CacheMap struct {
		Val       []byte
		CreatedAt time.Time
}


func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]CacheMap),
	}

	go c.ReapLoop(interval)

	return c
}


func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = CacheMap{
		Val: val,
		CreatedAt: time.Now().UTC(),
	};
}

func (c *Cache) Get(key string) ([]byte, bool)  {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.cache[key];

	return val.Val, ok
}


func (c *Cache) ReapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval);
	for range ticker.C {
		c.Reap(interval)
	}
} 


func (c *Cache) Reap(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock();
	
	for key, value := range c.cache{
		if value.CreatedAt.Add(interval).After(time.Now()){
			delete(c.cache, key)
		}
	}
}