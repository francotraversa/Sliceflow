package domain

import "github.com/francotraversa/Sliceflow/internal/types"

type OwnerUseCase interface {
	GetAllUsers() (*[]types.User, error)
}
type OwnerRepository interface {
	GetAllUsers() (*[]types.User, error)
}
