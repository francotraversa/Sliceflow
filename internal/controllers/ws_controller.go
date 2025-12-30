package controller

import (
	"context"
	"net/http"
	"time"

	redis "github.com/francotraversa/Sliceflow/internal/cache"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler mantiene la conexión abierta y envía updates
func WebSocketHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ctx := context.Background()
	pubsub := redis.RedisClient.Subscribe(ctx, "dashboard_updates")
	defer pubsub.Close()

	redisChannel := pubsub.Channel()

	for {
		select {
		case msg := <-redisChannel:
			if err := ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				c.Logger().Error("Cliente desconectado o error de escritura:", err)
				return nil
			}

		case <-time.After(30 * time.Second):
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return nil
			}
		}
	}
}
