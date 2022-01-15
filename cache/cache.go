package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gitlab.com/www.monitordomercado.com.br/mm/src/configuration"
)

var ctx = context.Background()
var initialized bool

var (
	// CacheDB - Redis cache connection
	CacheDB *redis.Client
)

// Init - Start the redis cache connection
func Init() {
	if !initialized {
		CacheDB = redis.NewClient(&redis.Options{
			Addr:     configuration.CFGs.SITE_CACHE_ADDR, // ex localhost:6379
			Password: "",
			DB:       1,
		})

		initialized = true
	}
}

// GetItem - Get item from redis cache
func GetItem(key string) (string, error) {
	return CacheDB.Get(ctx, key).Result()
}

// SetItem - Set item in redis cache
func SetItem(key string, value string) error {
	return CacheDB.Set(ctx, key, value, 0).Err()
}
