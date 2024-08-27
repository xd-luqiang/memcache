package memcache

import (
	"testing"
	"time"
)

const (
	maxTestEntries = 10
)

func TestNewCache(t *testing.T) {
	cache, reporter := NewCacheWithReport(WithCacheName("test"), WithShardings(4), WithLRU(4))
	v, found := cache.Get("1")
	if found {
		t.Logf("cache.Get(1) %v found", v)
	} else {
		t.Logf("cache.Get(1) not found")
	}
	cache.Set("1", 1, time.Second)
	v, found = cache.Get("1")
	if found {
		t.Logf("cache.Get(1) %v found", v)
	} else {
		t.Logf("cache.Get(1) not found")
	}
	time.Sleep(time.Second)
	v, found = cache.Get("1")
	if found {
		t.Logf("cache.Get(1) %v found", v)
	} else {
		t.Logf("cache.Get(1) not found")
	}
	t.Logf("cache.CacheName() %v", reporter.CacheName())
	t.Logf("cache.CacheType() %v", reporter.CacheType())
	t.Logf("cache.CacheShardings() %v", reporter.CacheShardings())
	t.Logf("cache.Size() %v", cache.Size())
}
