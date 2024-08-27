package memcache

import "time"

var (
	mapInitialCap   = 64
	sliceInitialCap = 64
)

func hash(key string) int {
	hash := 1469598103934665603

	for _, r := range key {
		hash = (hash << 5) - hash + int(r&0xffff)
		hash *= 1099511628211
	}

	return hash
}

func now() int64 {
	return time.Now().UnixNano()
}

// SetMapInitialCap sets the initial capacity of map.
func SetMapInitialCap(initialCap int) {
	if initialCap > 0 {
		mapInitialCap = initialCap
	}
}

// SetSliceInitialCap sets the initial capacity of slice.
func SetSliceInitialCap(initialCap int) {
	if initialCap > 0 {
		sliceInitialCap = initialCap
	}
}
