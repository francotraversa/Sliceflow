package domain

import "github.com/francotraversa/Sliceflow/internal/types"

// AuthUseCase define la lógica de negocio de autenticación
type AuthUseCase interface {
	Login(userCreds types.UserLoginCreds) (*types.TokenResponse, error)
}

// AuthRepository define las operaciones de base de datos para auth
type AuthRepository interface {
	CheckUser(userCreds types.UserLoginCreds) (bool, error)
	CheckPassword(userCreds types.UserLoginCreds) (bool, error)
	GetUser(username string) (*types.User, error)
}
