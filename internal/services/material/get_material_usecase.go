package services

import (
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllMaterialsUseCase(filter types.MaterialFilter) (*[]types.Material, error) {
	db := storage.DatabaseInstance{}.Instance()
	var materials []types.Material

	query := db.Model(&types.Material{})

	// Filtro por Nombre (Búsqueda parcial)
	if filter.Name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(filter.Name)+"%")
	}
	if filter.Type != "" {
		query = query.Where("LOWER(type) LIKE ?", "%"+strings.ToLower(filter.Type)+"%")
	}
	if filter.Brand != "" {
		query = query.Where("LOWER(brand) LIKE ?", "%"+strings.ToLower(filter.Brand)+"%")
	}

	if err := query.Find(&materials).Error; err != nil {
		return nil, err
	}
	return &materials, nil
}
