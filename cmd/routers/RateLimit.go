package routers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// NewRateLimiter crea un rate limiter configurable.
// rate = requests permitidos por segundo
// burst = pico máximo acumulable antes de empezar a rechazar
func NewRateLimiter(ratePerSec rate.Limit, burst int) echo.MiddlewareFunc {
	return echoMiddleware.RateLimiterWithConfig(echoMiddleware.RateLimiterConfig{
		Store: echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
			echoMiddleware.RateLimiterMemoryStoreConfig{
				Rate:      ratePerSec,
				Burst:     burst,
				ExpiresIn: 3 * time.Minute,
			},
		),
		ErrorHandler: func(ctx echo.Context, err error) error {
			return ctx.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "too many requests, slow down",
			})
		},
		DenyHandler: func(ctx echo.Context, identifier string, err error) error {
			return ctx.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "rate limit exceeded",
			})
		},
	})
}

func RateLimitMiddleware() echo.MiddlewareFunc {
	return NewRateLimiter(20, 50)
}

func TimeoutMiddleware() echo.MiddlewareFunc {
	return echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Timeout:      10 * time.Second,
		ErrorMessage: `{"error":"request timeout"}`,
	})
}
