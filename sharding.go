package memcache

import (
	"math/bits"
	"time"
)

type shardingCache struct {
	*config
	caches []Cache
}

func newShardingCache(conf *config, newCache func(conf *config) Cache) Cache {
	if conf.shardings <= 0 {
		panic("cachego: shardings must be > 0.")
	}

	if bits.OnesCount(uint(conf.shardings)) > 1 {
		panic("cachego: shardings must be the pow of 2 (such as 64).")
	}

	caches := make([]Cache, 0, conf.shardings)
	for i := 0; i < conf.shardings; i++ {
		caches = append(caches, newCache(conf))
	}

	cache := &shardingCache{
		config: conf,
		caches: caches,
	}

	return cache
}

func (sc *shardingCache) cacheOf(key string) Cache {
	hash := sc.hash(key)
	mask := len(sc.caches) - 1

	return sc.caches[hash&mask]
}

// Get gets the value of key from cache and returns value if found.
func (sc *shardingCache) Get(key string) (value interface{}, found bool) {
	return sc.cacheOf(key).Get(key)
}

// Set sets key and value to cache with ttl and returns evicted value if exists and unexpired.
// See Cache interface.
func (sc *shardingCache) Set(key string, value interface{}, ttl time.Duration) (oldValue interface{}) {
	return sc.cacheOf(key).Set(key, value, ttl)
}

// Remove removes key and returns the removed value of key.
// See Cache interface.
func (sc *shardingCache) Remove(key string) (removedValue interface{}) {
	return sc.cacheOf(key).Remove(key)
}

// Size returns the count of keys in cache.
// See Cache interface.
func (sc *shardingCache) Size() (size int) {
	for _, cache := range sc.caches {
		size += cache.Size()
	}

	return size
}

// GC cleans the expired keys in cache and returns the exact count cleaned.
// See Cache interface.
func (sc *shardingCache) GC() (cleans int) {
	for _, cache := range sc.caches {
		cleans += cache.GC()
	}

	return cleans
}

// Reset resets cache to initial status which is like a new cache.
// See Cache interface.
func (sc *shardingCache) Reset() {
	for _, cache := range sc.caches {
		cache.Reset()
	}
}

// Load loads a value by load function and sets it to cache.
// Returns an error if load failed.
// func (sc *shardingCache) Load(key string, ttl time.Duration, load func() (value interface{}, err error)) (value interface{}, err error) {
// 	return sc.cacheOf(key).Load(key, ttl, load)
// }
