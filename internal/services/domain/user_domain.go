package domain

import "github.com/francotraversa/Sliceflow/internal/types"

type UserRepository interface {
	GetUserByID(id uint) (*types.User, error)
	GetUserByUsername(username string) (*types.User, error)
	CreateUser(user *types.User) error
	UpdateUser(id uint, user *types.UserUpdateCreds, companyID uint) error
	DeleteUser(id uint, companyID uint) error
	GetUsers(companyID uint) ([]types.User, error)
}

type UserUseCase interface {
	GetUserByID(id uint) (*types.User, error)
	GetUserByUsername(username string) (*types.User, error)
	CreateUser(user *types.UserCreateCreds, companyID uint) error
	UpdateUser(id uint, user *types.UserUpdateCreds, companyID uint) error
	DeleteUser(id uint, companyID uint) error
	GetUsers(companyID uint) ([]types.User, error)
}
