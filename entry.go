package memcache

import "time"

type entry struct {
	key        string
	value      interface{}
	expiration int64 // Time in nanosecond, valid util 2262 year (enough, uh?)
	now        func() int64
}

func newEntry(key string, value interface{}, ttl time.Duration, now func() int64) *entry {
	e := &entry{
		now: now,
	}

	e.setup(key, value, ttl)
	return e
}

func (e *entry) setup(key string, value interface{}, ttl time.Duration) {
	e.key = key
	e.value = value
	e.expiration = 0

	if ttl > 0 {
		e.expiration = e.now() + ttl.Nanoseconds()
	}
}

func (e *entry) expired(now int64) bool {
	if now > 0 {
		return e.expiration > 0 && e.expiration < now
	}

	return e.expiration > 0 && e.expiration < e.now()
}
