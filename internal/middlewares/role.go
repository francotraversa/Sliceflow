package middlewares

import (
	"net/http"

	"github.com/francotraversa/Sliceflow/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tok, ok := c.Get("user").(*jwt.Token)
			if !ok || tok == nil || !tok.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
			}

			claims, ok := tok.Claims.(*auth.JwtCustomClaims)
			if !ok || claims == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid claims")
			}

			if claims.Role != role {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}

			return next(c)
		}
	}
}
