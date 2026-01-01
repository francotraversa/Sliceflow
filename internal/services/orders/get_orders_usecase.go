package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllOrdersUseCase(filter types.OrderFilter) (*[]types.ProductionOrder, error) {
	db := storage.DatabaseInstance{}.Instance()
	cacheKey := fmt.Sprintf("orders:list:%s", filter.Status)
	var orders []types.ProductionOrder

	if services.GetCache(cacheKey, &orders) {
		return &orders, nil
	}

	query := db.Preload("Material").Preload("Machine")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.SortPriority {
		query = query.Order("priority ASC") // P1 primero
	} else {
		query = query.Order("created_at DESC") // Las más nuevas primero
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	services.SetCache(cacheKey, &orders)
	return &orders, nil
}
