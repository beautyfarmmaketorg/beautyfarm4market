package util

import (
	cache "github.com/UncleBig/goCache"
	"time"
)

var currentCache *cache.Cache

func init() {
	currentCache = cache.New(60*time.Minute, 30*time.Second)
}

func SetCache(key string, value interface{}, d time.Duration) {
	currentCache.Set(key, value, d)
}

func GetCache(key string) interface{} {
	if cacheValue, found := currentCache.Get(key); found {
		return cacheValue
	}
	return nil
}
