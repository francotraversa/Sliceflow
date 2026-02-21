package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisPass := os.Getenv("REDIS_PASSWORD")

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	fmt.Printf("🔌 Conectando a Redis en: %s ...\n", addr)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		fmt.Printf("⚠️  Error conecting Redis with Backend: %v\n", err)
	} else {
		fmt.Println("Conexion to Redis successfully!")
	}
}
