package memcache

import "time"

type (
	DeserializeFunc func(string, []byte) interface{}
	LoadFunc        func([]string, DeserializeFunc) ([]interface{}, error)
)

type config struct {
	cacheName    string
	cacheType    CacheType
	shardings    int
	singleflight bool
	expireTime   time.Duration
	protectTime  time.Duration
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

	loadFunc LoadFunc
}

func newDefaultConfig() *config {
	return &config{
		cacheName:    "",
		cacheType:    standard,
		shardings:    0,
		singleflight: true,
		expireTime:   60 * time.Second,
		protectTime:  10 * time.Second,
		gcDuration:   10 * time.Minute,
		maxScans:     10000,
		maxEntries:   100000,
		now:          now,
		hash:         hash,
		recordMissed: true,
		recordHit:    true,
		recordGC:     true,
		recordLoad:   true,
		loadFunc:     nil,
	}
}
