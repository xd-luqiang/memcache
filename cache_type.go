package memcache

const (
	// standard cache is a simple cache with locked map.
	// It evicts entries randomly if cache size reaches to max entries.
	standard CacheType = "standard"

	// lru cache is a cache using lru to evict entries.
	// More details see https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU).
	lru CacheType = "lru"

	// lfu cache is a cache using lfu to evict entries.
	// More details see https://en.wikipedia.org/wiki/Cache_replacement_policies#Least-frequently_used_(LFU).
	lfu CacheType = "lfu"
)

// CacheType is the type of cache.
type CacheType string

// String returns the cache type in string form.
func (ct CacheType) String() string {
	return string(ct)
}

// IsStandard returns if cache type is standard.
func (ct CacheType) IsStandard() bool {
	return ct == standard
}

// IsLRU returns if cache type is lru.
func (ct CacheType) IsLRU() bool {
	return ct == lru
}

// IsLFU returns if cache type is lfu.
func (ct CacheType) IsLFU() bool {
	return ct == lfu
}
