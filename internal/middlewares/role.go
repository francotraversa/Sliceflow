package middlewares

import (
	"net/http"

	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := GetClaimsFromContext(c)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, types.Error{Error: "failed to parse custom claims"})
			}

			if claims.Role != role && claims.Role != "owner" {
				return c.JSON(http.StatusForbidden, types.Error{Error: "permission denied: only admins can access the user list"})
			}

			return next(c)
		}
	}
}
