package services

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/francotraversa/Sliceflow/internal/auth"
	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/francotraversa/Sliceflow/internal/utils"
)

func AuthUseCase(userCread types.UserLoginCreds) (*types.TokenResponse, error) {
	if userCread.Username == "" || userCread.Password == "" {
		return nil, fmt.Errorf("Insufficient parameters")
	}

	User := storage.FindUserByUsername(strings.ToLower(strings.TrimSpace(userCread.Username)))
	if User == nil {
		return nil, fmt.Errorf("Invalid username or password")
	}

	if User.Status == "disabled" {
		return nil, fmt.Errorf("this account is disabled, please contact support")
	}

	err := utils.CheckPassword(User.Password, userCread.Password)
	if err != nil {
		return nil, fmt.Errorf("Invalid username or password")
	}

	ttl, err := strconv.Atoi(os.Getenv("TTL"))
	if err != nil {
		return nil, fmt.Errorf("Invalid TTL value")
	}

	token, err := auth.GenerateToken(User.ID, User.Role, os.Getenv("JWT_SECRET"), time.Duration(ttl)*time.Minute)
	if err != nil {
		return nil, err
	}

	return token, nil
}
