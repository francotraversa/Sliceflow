package middlewares

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetClaimsFromContext(c echo.Context) (*types.JwtCustomClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok || token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*types.JwtCustomClaims)
	if !ok || claims == nil {
		return nil, fmt.Errorf("invalid jwt claims")
	}

	return claims, nil
}
