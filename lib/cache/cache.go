package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var cache = ttlcache.New(
	ttlcache.WithTTL[string, string](0),
)

func InitializeNewCache(accessTokenDuration time.Duration) *ttlcache.Cache[string, string] {
	return ttlcache.New(
		ttlcache.WithTTL[string, string](time.Minute * accessTokenDuration),
	)
}

func SetCache(key string, value string) {
	cache.Set(key, value, ttlcache.DefaultTTL)
}

func GetCache(key string) *ttlcache.Item[string, string] {
	return cache.Get(key)
}

func DeleteCache(key string) {
	cache.Delete(key)
}
