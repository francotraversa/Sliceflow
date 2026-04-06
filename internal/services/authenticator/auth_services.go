package services

import "github.com/francotraversa/Sliceflow/internal/services/domain"

type authUseCase struct {
	repo domain.AuthRepository
}

func NewAuthUseCase(repo domain.AuthRepository) domain.AuthUseCase {
	return &authUseCase{repo: repo}
}
