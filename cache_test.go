package memcache

import (
	"fmt"
	"testing"
	"time"
)

func testDeserializeFunc(key string, data []byte) interface{} {
	return string(data) + key
}

func testLoadfunc(keys []string, deserializeF DeserializeFunc) ([]interface{}, error) {
	fmt.Println(keys)
	res := make([]interface{}, len(keys))
	for i, key := range keys {
		res[i] = key + "value"
		if deserializeF != nil {
			res[i] = deserializeF(key, []byte(res[i].(string)))
		}
	}
	return res, nil
}

func TestNewCache(t *testing.T) {
	cache, reporter := NewCacheWithReport(WithCacheName("test"), WithShardings(4), WithLRU(4), WithGC(0), WithExpire(5*time.Second), WithProtect(time.Second), withLoadFunc(testLoadfunc))
	cancell := RunGCTask(cache, 500*time.Millisecond)
	keys := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	vals := []interface{}{1, 2, 3, nil, 5, nil, 7, 8}
	cache.MSet(keys, vals)

	qkeys := []string{"1", "2", "3", "4", "5", "6", "7", "8", "0", "10"}
	values, founds := cache.MGet(qkeys, testDeserializeFunc)
	t.Logf("cache.Get %v", values)
	t.Logf("cache.Get %v", founds)
	// t.Logf("cache.Size() %v", cache.Size())

	time.Sleep(1*time.Second + 100*time.Millisecond)
	values, founds = cache.MGet(qkeys, testDeserializeFunc)
	t.Logf("cache.Get %v", values)
	t.Logf("cache.Get %v", founds)
	// t.Logf("cache.Size() %v", cache.Size())

	// time.Sleep(4 * time.Second)
	// t.Logf("cache.Size() %v", cache.Size())

	cancell()
	// t.Logf("cache.Size() %v", cache.Size())

	// // t.Logf("cache.CacheName() %v", reporter.CacheName())
	// // t.Logf("cache.CacheType() %v", reporter.CacheType())
	// // t.Logf("cache.CacheShardings() %v", reporter.CacheShardings())
	// t.Logf("cache.Size() %v", cache.Size())
	t.Logf("cache.MissedRate() %v", reporter.MissedRate())
}
