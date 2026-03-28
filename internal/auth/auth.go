package auth

import (
	"time"

	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, companyID uint, user string, role string, secret string, ttl time.Duration) (*types.TokenResponse, error) {
	exp := time.Now().Add(ttl).Unix()

	claims := &types.JwtCustomClaims{
		UserId:    userID,
		Role:      role,
		User:      user,
		CompanyId: companyID,
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
