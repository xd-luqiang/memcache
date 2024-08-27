package memcache

import "time"

type config struct {
	cacheName    string
	cacheType    CacheType
	shardings    int
	singleflight bool
	gcDuration   time.Duration

	maxScans   int
	maxEntries int

	now  func() int64
	hash func(key string) int

	recordMissed bool
	recordHit    bool
	recordGC     bool
	recordLoad   bool

	reportMissed func(reporter *Reporter, key string)
	reportHit    func(reporter *Reporter, key string, value interface{})
	reportGC     func(reporter *Reporter, cost time.Duration, cleans int)
	reportLoad   func(reporter *Reporter, key string, value interface{}, ttl time.Duration, err error)
}

func newDefaultConfig() *config {
	return &config{
		cacheName:    "",
		cacheType:    standard,
		shardings:    0,
		singleflight: true,
		gcDuration:   10 * time.Minute,
		maxScans:     10000,
		maxEntries:   100000,
		now:          now,
		hash:         hash,
		recordMissed: true,
		recordHit:    true,
		recordGC:     true,
		recordLoad:   true,
	}
}
