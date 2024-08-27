package memcache

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache, reporter := NewCacheWithReport(WithCacheName("test"), WithShardings(4), WithLRU(4), WithGC(0), WithExpire(5*time.Second), WithProtect(time.Second))
	cancell := RunGCTask(cache, 500*time.Millisecond)
	keys := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	vals := []interface{}{1, 2, 3, nil, 5, nil, 7, 8}
	cache.MSet(keys, vals)
	t.Logf("cache.Size() %v", cache.Size())

	time.Sleep(1*time.Second + 100*time.Millisecond)
	t.Logf("cache.Size() %v", cache.Size())

	time.Sleep(4 * time.Second)
	t.Logf("cache.Size() %v", cache.Size())

	cancell()
	t.Logf("cache.Size() %v", cache.Size())

	// qkeys := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	// values, founds := cache.MGet(qkeys)
	// t.Logf("cache.Get %v", values)
	// t.Logf("cache.Get %v", founds)
	// // t.Logf("cache.CacheName() %v", reporter.CacheName())
	// // t.Logf("cache.CacheType() %v", reporter.CacheType())
	// // t.Logf("cache.CacheShardings() %v", reporter.CacheShardings())
	// t.Logf("cache.Size() %v", cache.Size())
	t.Logf("cache.MissedRate() %v", reporter.MissedRate())
}
