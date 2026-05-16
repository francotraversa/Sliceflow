package services

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/francotraversa/Sliceflow/internal/types"
)

type OwnerService struct {
	ownerRepo domain.OwnerRepository
}

func (s *OwnerService) GetAllUsers() (*[]types.User, error) {
	users, err := s.ownerRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}

	return users, nil
}

func NewOwnerService(ownerRepo domain.OwnerRepository) *OwnerService {
	return &OwnerService{ownerRepo: ownerRepo}
}
