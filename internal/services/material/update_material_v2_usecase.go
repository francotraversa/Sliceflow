package services

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *materialService) UpdateMaterial(id uint, dto types.UpdateMaterialDTO, companyID uint) error {
	material, err := s.repo.GetByID(id, companyID)
	if err != nil {
		return fmt.Errorf("The Material was not found")
	}
	if dto.Name != "" {
		material.Name = dto.Name
	}
	if dto.Type != "" {
		material.Type = dto.Type
	}
	if dto.Brand != "" {
		material.Brand = dto.Brand
	}
	if dto.Description != "" {
		material.Description = dto.Description
	}

	return s.repo.Update(id, dto, companyID)
}
