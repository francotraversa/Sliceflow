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
// @Summary      WebSocket connection for Dashboard
// @Description  Establishes a persistent WebSocket connection. Listens to Redis events (channel "dashboard_updates") and sends them to the client in real-time.
// @Description  <br> **Auth:** Enviar el JWT como query param: `?token=<jwt>` (los WebSockets no soportan headers de autenticación).
// @Description  <br> **Expected events:** `REFRESH_ORDERS`, `MACHINE_UPDATED`, etc.
// @Description  <br> **Note:** Swagger UI does not natively support WebSocket testing. Use Bruno, Postman or PieSocket.
// @Tags         Dashboard Realtime
// @Accept       json
// @Produce      json
// @Param        token query string true "JWT token"
// @Success      101  {string}  string  "Switching Protocols (Connection Established)"
// @Failure      401  {string}  string  "Token inválido o ausente"
// @Failure      400  {string}  string  "Protocol upgrade failed"
// @Failure      500  {string}  string  "Error interno del servidor"
// @Router       /hornero/authed/ws/dashboard [get]
func WebSocketHandler(c echo.Context) error {
	// El middleware echojwt ya validó el token (vía header o ?token= query param)
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
				c.Logger().Error("Client disconnected or write error:", err)
				return nil
			}

		case <-time.After(30 * time.Second):
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return nil
			}
		}
	}
}
