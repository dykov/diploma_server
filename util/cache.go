package util

import (
	gocache "github.com/pmylund/go-cache"
	"net/http"
)

var cache = gocache.New(cacheDefaultExpiration, cacheCleanupInterval)

func Cache(cacheKey string) (interface{}, bool) {
	if cacheValue, foundKey := cache.Get(cacheKey); foundKey {
		return cacheValue, true
	}
	return nil, false
}

func SetCache(cacheKey string, cacheValue interface{}) {
	cache.Set(cacheKey, cacheValue, gocache.DefaultExpiration)
}

func DeleteAllCache(w http.ResponseWriter, r *http.Request) {
	cache.Flush()
	w.WriteHeader(http.StatusNoContent)
}

func DeleteCache(cacheKey string) {
	cache.Delete(cacheKey)
}
