package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu       sync.Mutex
	cacheMap map[string]cacheEntry
}

func (currCache *Cache) Add(key string, val []byte) {
	fmt.Printf("Saving key %s\n", key)
	currCache.mu.Lock()
	defer currCache.mu.Unlock()
	newEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	currCache.cacheMap[key] = newEntry
}

func (currCache *Cache) Get(key string) ([]byte, bool) {
	fmt.Printf("Retrieving key %s\n", key)
	currCache.mu.Lock()
	defer currCache.mu.Unlock()
	entry, ok := currCache.cacheMap[key]
	return entry.val, ok
}

func (currCache *Cache) reap(interval time.Duration) {
	fmt.Println("Reaping")
	currCache.mu.Lock()
	defer currCache.mu.Unlock()
	for key, val := range currCache.cacheMap {
		elapsed := time.Since(val.createdAt)
		if elapsed > interval {
			delete(currCache.cacheMap, key)
		}
	}
}

func (currCache *Cache) readLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		currCache.reap(interval)
	}
}

func NewCache(interval time.Duration) *Cache {
	returnCache := &Cache{
		cacheMap: make(map[string]cacheEntry),
	}
	go returnCache.readLoop(interval)
	return returnCache
}
