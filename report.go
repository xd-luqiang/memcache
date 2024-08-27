package memcache

import (
	"sync/atomic"
	"time"
)

// Reporter stores some values for reporting.
type Reporter struct {
	conf  *config
	cache Cache

	missedCount uint64
	hitCount    uint64
	gcCount     uint64
	loadCount   uint64
}

func (r *Reporter) increaseMissedCount() {
	atomic.AddUint64(&r.missedCount, 1)
}

func (r *Reporter) increaseHitCount() {
	atomic.AddUint64(&r.hitCount, 1)
}

func (r *Reporter) increaseGCCount() {
	atomic.AddUint64(&r.gcCount, 1)
}

func (r *Reporter) increaseLoadCount() {
	atomic.AddUint64(&r.loadCount, 1)
}

// CacheName returns the name of cache.
// You can use WithCacheName to set cache's name.
func (r *Reporter) CacheName() string {
	return r.conf.cacheName
}

// CacheType returns the type of cache.
// See CacheType.
func (r *Reporter) CacheType() CacheType {
	return r.conf.cacheType
}

// CacheShardings returns the shardings of cache.
// You can use WithShardings to set cache's shardings.
// Zero shardings means cache is non-sharding.
func (r *Reporter) CacheShardings() int {
	return r.conf.shardings
}

// CacheGC returns the gc duration of cache.
// You can use WithGC to set cache's gc duration.
// Zero duration means cache disables gc.
func (r *Reporter) CacheGC() time.Duration {
	return r.conf.gcDuration
}

// CacheSize returns the size of cache.
func (r *Reporter) CacheSize() int {
	return r.cache.Size()
}

// CountMissed returns the missed count.
func (r *Reporter) CountMissed() uint64 {
	return atomic.LoadUint64(&r.missedCount)
}

// CountHit returns the hit count.
func (r *Reporter) CountHit() uint64 {
	return atomic.LoadUint64(&r.hitCount)
}

// CountGC returns the gc count.
func (r *Reporter) CountGC() uint64 {
	return atomic.LoadUint64(&r.gcCount)
}

// CountLoad returns the load count.
func (r *Reporter) CountLoad() uint64 {
	return atomic.LoadUint64(&r.loadCount)
}

// MissedRate returns the missed rate.
func (r *Reporter) MissedRate() float64 {
	hit := r.CountHit()
	missed := r.CountMissed()

	total := hit + missed
	if total <= 0 {
		return 0.0
	}

	return float64(missed) / float64(total)
}

// HitRate returns the hit rate.
func (r *Reporter) HitRate() float64 {
	hit := r.CountHit()
	missed := r.CountMissed()

	total := hit + missed
	if total <= 0 {
		return 0.0
	}

	return float64(hit) / float64(total)
}

type reportableCache struct {
	*config
	*Reporter
}

func report(conf *config, cache Cache) (Cache, *Reporter) {
	reporter := &Reporter{
		conf:        conf,
		cache:       cache,
		hitCount:    0,
		missedCount: 0,
		gcCount:     0,
		loadCount:   0,
	}

	cache = &reportableCache{
		config:   conf,
		Reporter: reporter,
	}

	return cache, reporter
}

// Get gets the value of key from cache and returns value if found.
func (rc *reportableCache) Get(key string, deserializeF DeserializeFunc) (value interface{}, found bool) {
	value, found = rc.cache.Get(key, deserializeF)

	if found {
		if rc.recordHit {
			rc.increaseHitCount()
		}

		if rc.reportHit != nil {
			rc.reportHit(rc.Reporter, key, value)
		}
	} else {
		if rc.recordMissed {
			rc.increaseMissedCount()
		}

		if rc.reportMissed != nil {
			rc.reportMissed(rc.Reporter, key)
		}
	}

	return value, found
}

func (rc *reportableCache) MGet(keys []string, deserializeF DeserializeFunc) (values []interface{}, founds []bool) {
	values, founds = rc.cache.MGet(keys, deserializeF)
	for i, found := range founds {
		if found {
			if rc.recordHit {
				rc.increaseHitCount()
			}

			if rc.reportHit != nil {
				rc.reportHit(rc.Reporter, keys[i], values[i])
			}
		} else {
			if rc.recordMissed {
				rc.increaseMissedCount()
			}

			if rc.reportMissed != nil {
				rc.reportMissed(rc.Reporter, keys[i])
			}
		}
	}

	return values, founds
}

// Set sets key and value to cache with ttl and returns evicted value if exists and unexpired.
// See Cache interface.
func (rc *reportableCache) Set(key string, value interface{}, ttl ...time.Duration) (evictedValue interface{}) {
	return rc.cache.Set(key, value, ttl...)
}

func (rc *reportableCache) MSet(keys []string, values []interface{}, ttls ...time.Duration) (evictedValues []interface{}) {
	evictedValues = rc.cache.MSet(keys, values, ttls...)
	return evictedValues
}

// Remove removes key and returns the removed value of key.
// See Cache interface.
func (rc *reportableCache) Remove(key string) (removedValue interface{}) {
	return rc.cache.Remove(key)
}

// Size returns the count of keys in cache.
// See Cache interface.
func (rc *reportableCache) Size() (size int) {
	return rc.cache.Size()
}

// GC cleans the expired keys in cache and returns the exact count cleaned.
// See Cache interface.
func (rc *reportableCache) GC() (cleans int) {
	if rc.recordGC {
		rc.increaseGCCount()
	}

	if rc.reportGC == nil {
		return rc.cache.GC()
	}

	begin := rc.now()
	cleans = rc.cache.GC()
	end := rc.now()

	cost := time.Duration(end - begin)
	rc.reportGC(rc.Reporter, cost, cleans)

	return cleans
}

// Reset resets cache to initial status which is like a new cache.
// See Cache interface.
func (rc *reportableCache) Reset() {
	rc.cache.Reset()
}

// Load loads a key with ttl to cache and returns an error if failed.
// See Cache interface.
// func (rc *reportableCache) Load(key string, ttl time.Duration, load func() (value interface{}, err error)) (value interface{}, err error) {
// 	value, err = rc.cache.Load(key, ttl, load)

// 	if rc.recordLoad {
// 		rc.increaseLoadCount()
// 	}

// 	if rc.reportLoad != nil {
// 		rc.reportLoad(rc.Reporter, key, value, ttl, err)
// 	}

// 	return value, err
// }
