package services

import (
	"context"
	"fmt"

	redis "github.com/francotraversa/Sliceflow/internal/cache"
)

func PublishEvent(channel string, message string) {
	if redis.RedisClient == nil {
		return
	}

	// Enviamos el mensaje a Redis (Fire and Forget)
	err := redis.RedisClient.Publish(context.Background(), channel, message).Err()
	if err != nil {
		fmt.Printf("❌ Error publishing WS event: %v\n", err)
	} else {
		fmt.Printf("📢 WS event sent: %s -> %s\n", channel, message)
	}
}
