package memcache

import (
	"container/list"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type lruCache struct {
	*config

	elementMap  map[string]*list.Element
	elementList *list.List
	lock        sync.RWMutex

	// loader *loader
}

func newLRUCache(conf *config) Cache {
	if conf.maxEntries <= 0 {
		panic("cachego: lru cache must specify max entries")
	}

	cache := &lruCache{
		config:      conf,
		elementMap:  make(map[string]*list.Element, mapInitialCap),
		elementList: list.New(),
		// loader:      newLoader(conf.singleflight),
	}

	return cache
}

func (lc *lruCache) unwrap(element *list.Element) *entry {
	entry, ok := element.Value.(*entry)
	if !ok {
		panic("cachego: failed to unwrap lru element's value to entry")
	}

	return entry
}

func (lc *lruCache) evict() (evictedValue interface{}) {
	if element := lc.elementList.Back(); element != nil {
		return lc.removeElement(element)
	}

	return nil
}

func (lc *lruCache) get(key string) (value interface{}, found bool) {
	element, ok := lc.elementMap[key]
	if !ok {
		return nil, false
	}

	entry := lc.unwrap(element)
	if entry.expired(0) {
		return nil, false
	}
	lc.elementList.MoveToFront(element)
	if entry.value == nil {
		return nil, false
	} else {
		return *entry.value, true
	}
}

func fluctuate(a time.Duration) time.Duration {
	// 计算上下 20% 的浮动范围
	fluctuation := int64(float64(a) * 0.2)

	// 生成一个在 [-fluctuation, fluctuation] 范围内的随机浮动值
	delta := rand.Int63n(2*fluctuation+1) - fluctuation

	// 将浮动值添加到 a 的值上
	return a + time.Duration(delta)
}

func (lc *lruCache) set(key string, value interface{}, ttl ...time.Duration) (evictedValue interface{}) {
	curTtl := fluctuate(lc.expireTime)
	if len(ttl) > 0 {
		curTtl = ttl[0]
	}
	if value == nil {
		curTtl = lc.protectTime
	}
	fmt.Println(key, curTtl)
	element, ok := lc.elementMap[key]
	if ok {
		entry := lc.unwrap(element)
		entry.setup(key, &value, curTtl)

		lc.elementList.MoveToFront(element)
		return nil
	}

	if lc.maxEntries > 0 && lc.elementList.Len() >= lc.maxEntries {
		evictedValue = lc.evict()
	}

	element = lc.elementList.PushFront(newEntry(key, &value, curTtl, lc.now))
	lc.elementMap[key] = element

	return evictedValue
}

func (lc *lruCache) removeElement(element *list.Element) (removedValue interface{}) {
	entry := lc.unwrap(element)

	delete(lc.elementMap, entry.key)
	lc.elementList.Remove(element)

	return entry.value
}

func (lc *lruCache) remove(key string) (removedValue interface{}) {
	if element, ok := lc.elementMap[key]; ok {
		return lc.removeElement(element)
	}

	return nil
}

func (lc *lruCache) size() (size int) {
	return len(lc.elementMap)
}

func (lc *lruCache) gc() (cleans int) {
	now := lc.now()
	scans := 0

	for _, element := range lc.elementMap {
		scans++

		if entry := lc.unwrap(element); entry.expired(now) {
			lc.removeElement(element)
			cleans++
		}

		if lc.maxScans > 0 && scans >= lc.maxScans {
			break
		}
	}

	return cleans
}

func (lc *lruCache) reset() {
	lc.elementMap = make(map[string]*list.Element, mapInitialCap)
	lc.elementList = list.New()

	// lc.loader.Reset()
}

// Get gets the value of key from cache and returns value if found.
// See Cache interface.
func (lc *lruCache) Get(key string) (value interface{}, found bool) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	return lc.get(key)
}

func (lc *lruCache) MGet(keys []string) (values []interface{}, founds []bool) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	for _, key := range keys {
		value, found := lc.get(key)
		values = append(values, value)
		founds = append(founds, found)
	}
	return values, founds
}

// Set sets key and value to cache with ttl and returns evicted value if exists and unexpired.
// See Cache interface.
func (lc *lruCache) Set(key string, value interface{}, ttl ...time.Duration) (evictedValue interface{}) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	return lc.set(key, value, ttl...)
}

func (lc *lruCache) MSet(keys []string, values []interface{}, ttls ...time.Duration) (evictedValues []interface{}) {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	if len(keys) != len(values) {
		fmt.Printf("cachego: keys and values must have the same length, key: %v, value: %v\n", keys, values)
		return nil
	}
	for i := 0; i < len(keys); i++ {
		if len(ttls) > i {
			evictedValues = append(evictedValues, lc.set(keys[i], values[i], ttls[i]))
		} else {
			evictedValues = append(evictedValues, lc.set(keys[i], values[i]))
		}
	}
	return evictedValues
}

// Remove removes key and returns the removed value of key.
// See Cache interface.
func (lc *lruCache) Remove(key string) (removedValue interface{}) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	return lc.remove(key)
}

// Size returns the count of keys in cache.
// See Cache interface.
func (lc *lruCache) Size() (size int) {
	lc.lock.RLock()
	defer lc.lock.RUnlock()

	return lc.size()
}

// GC cleans the expired keys in cache and returns the exact count cleaned.
// See Cache interface.
func (lc *lruCache) GC() (cleans int) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	return lc.gc()
}

// Reset resets cache to initial status which is like a new cache.
// See Cache interface.
func (lc *lruCache) Reset() {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	lc.reset()
}

// Load loads a value by load function and sets it to cache.
// Returns an error if load failed.
// func (lc *lruCache) Load(key string, ttl time.Duration, load func() (value interface{}, err error)) (value interface{}, err error) {
// 	value, err = lc.loader.Load(key, ttl, load)
// 	if err != nil {
// 		return value, err
// 	}

// 	lc.Set(key, value, ttl)
// 	return value, nil
// }
