package controller

import (
	"context"
	"net/http"
	"time"

	redis "github.com/francotraversa/Sliceflow/internal/infra/cache"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler godoc
// @Summary      Conexión WebSocket para Dashboard
// @Description  Establece una conexión WebSocket persistente. Escucha eventos de Redis (canal "dashboard_updates") y los envía al cliente en tiempo real.
// @Description  <br> **Eventos esperados:** `REFRESH_ORDERS`, `MACHINE_UPDATED`, etc.
// @Description  <br> **Nota:** Swagger UI no soporta probar WebSockets nativamente. Usar Bruno, Postman o PieSocket.
// @Tags         Dashboard Realtime
// @Accept       json
// @Produce      json
// @Success      101  {string}  string  "Switching Protocols (Conexión Establecida)"
// @Failure      400  {string}  string  "Error al actualizar protocolo (Upgrade failed)"
// @Failure      500  {string}  string  "Error interno del servidor"
// @Router       /hornero/ws/dashboard [get]
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
