package middlewares

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	UserId uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func IdFromContext(c echo.Context) (uint, error) {
	tok, ok := c.Get("user").(*jwt.Token)
	if !ok || tok == nil || !tok.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := tok.Claims.(*JwtCustomClaims)
	if !ok || claims == nil {
		return 0, fmt.Errorf("invalid jwt claims")
	}

	return claims.UserId, nil
}
