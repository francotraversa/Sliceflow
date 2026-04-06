package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/francotraversa/Sliceflow/internal/auth"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func (uc *authUseCase) Login(userCreds types.UserLoginCreds) (*types.TokenResponse, error) {
	if userCreds.Username == "" || userCreds.Password == "" {
		return nil, fmt.Errorf("Insufficient parameters")
	}
	userExists, err := uc.repo.CheckUser(userCreds)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, fmt.Errorf("Invalid credentials")
	}

	passwordMatch, err := uc.repo.CheckPassword(userCreds)
	if err != nil {
		return nil, err
	}
	if !passwordMatch {
		return nil, fmt.Errorf("Invalid credentials")
	}

	ttl, err := strconv.Atoi(os.Getenv("TTL"))
	if err != nil {
		return nil, fmt.Errorf("Invalid TTL value")
	}

	user, err := uc.repo.GetUser(userCreds.Username)
	if err != nil {
		return nil, err
	}
	token, err := auth.GenerateToken(user.IdUser, user.IdCompany, user.Username, user.Role, os.Getenv("JWT_SECRET"), time.Duration(ttl)*time.Minute)
	if err != nil {
		return nil, err
	}

	return token, nil
}
