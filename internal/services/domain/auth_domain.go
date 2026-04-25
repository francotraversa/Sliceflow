package domain

import "github.com/francotraversa/Sliceflow/internal/types"

type AuthRepository interface {
	Login(userCreds types.UserLoginCreds) (*types.TokenResponse, error)
	CheckUser(userCreds types.UserLoginCreds) (bool, error)
	CheckPassword(userCreds types.UserLoginCreds) (bool, error)
	GetUser(username string) (*types.User, error)
}

type AuthUseCase interface {
	Login(userCreds types.UserLoginCreds) (*types.TokenResponse, error)
}
