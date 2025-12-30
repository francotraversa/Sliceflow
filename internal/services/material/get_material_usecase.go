package services

import (
	"fmt"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllMaterialsUseCase(filter types.MaterialFilter) (*[]types.Material, error) {
	db := storage.DatabaseInstance{}.Instance()
	cacheKey := fmt.Sprintf("materials:list:%s:%s", filter.Name, filter.Type)
	var materials []types.Material

	if services.GetCache(cacheKey, &materials) {
		return &materials, nil
	}

	query := db.Model(&types.Material{})

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
	services.SetCache(cacheKey, &materials)
	return &materials, nil
}
