package services

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *materialService) GetMaterialByID(id uint, companyID uint) (*types.Material, error) {
	material, err := s.repo.GetByID(id, companyID)
	if err != nil {
		return nil, fmt.Errorf("The Material was not found")
	}
	return material, nil
}

func (s *materialService) GetMaterials(filter types.MaterialFilter, companyID uint) ([]types.Material, error) {
	materials, err := s.repo.GetMaterials(filter, companyID)
	if err != nil {
		return nil, fmt.Errorf("The Materials were not found")
	}
	return materials, nil
}
