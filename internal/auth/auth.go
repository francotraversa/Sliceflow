package auth

import (
	"time"

	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	UserId uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role, secret string, ttl time.Duration) (*types.TokenResponse, error) {
	exp := time.Now().Add(ttl).Unix()

	claims := &JwtCustomClaims{
		UserId: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &types.TokenResponse{
		Token:   token,
		Expires: exp,
	}, nil
}
