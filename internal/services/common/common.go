package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/francotraversa/Sliceflow/internal/infra/cache"
)

var ctx = context.Background()

const DefaultTTL = 10 * time.Minute

func GetCache(key string, dest interface{}) bool {
	if redis.RedisClient == nil {
		return false
	}

	val, err := redis.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		fmt.Printf("❌ Error deserializing key %s: %v\n", key, err)
		return false
	}
	return true
}

func SetCache(key string, data interface{}) {
	if redis.RedisClient == nil {
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("❌ Error deserializing key %s: %v\n", key, err)
		return
	}

	redis.RedisClient.Set(ctx, key, bytes, DefaultTTL)
}

func InvalidateCache(pattern string) {
	if redis.RedisClient == nil {
		return
	}

	keys, err := redis.RedisClient.Keys(ctx, pattern).Result()
	if err == nil && len(keys) > 0 {
		redis.RedisClient.Del(ctx, keys...)
		fmt.Printf("🧹 Cache invalidated: %s\n", pattern)
	}
}
