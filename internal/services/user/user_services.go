package services

import "github.com/francotraversa/Sliceflow/internal/services/domain"

type UserServices struct {
	userRepo domain.UserRepository
}

func NewUserServices(userRepo domain.UserRepository) domain.UserUseCase {
	return &UserServices{userRepo: userRepo}
}
