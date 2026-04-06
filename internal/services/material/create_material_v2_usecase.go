package services

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *materialService) CreateMaterial(material types.CreateMaterialDTO, companyID uint) error {
	if material.Name == "" || material.Type == "" {
		return fmt.Errorf("Name and Type are required")
	}
	newMaterial := types.Material{
		Name:        material.Name,
		Type:        material.Type,
		Description: material.Description,
		Brand:       material.Brand,
		IdCompany:   companyID,
	}
	return s.repo.Create(&newMaterial)
}
