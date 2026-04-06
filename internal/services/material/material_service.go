package services

import (
	"github.com/francotraversa/Sliceflow/internal/services/domain"
)

type materialService struct {
	repo domain.MaterialRepository
}

func NewMaterialService(repo domain.MaterialRepository) domain.MaterialUseCase {
	return &materialService{repo: repo}
}
