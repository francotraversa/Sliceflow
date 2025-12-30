package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(host string, port string, password string) {
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("🔌 Conectando a Redis en: %s ...\n", addr)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // DB por defecto
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		fmt.Printf("⚠️  Error conecting Redis with Backend: %v\n", err)
	} else {
		fmt.Println("🚀 Conexion to Redis successfully!")
	}
}
