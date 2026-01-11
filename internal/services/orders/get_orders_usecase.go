package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllOrdersUseCase(filter types.OrderFilter) (*[]types.ProductionOrder, error) {
	db := storage.DatabaseInstance{}.Instance()

	cacheKey := fmt.Sprintf("orders:list:st_%s:id_%d:sort_%v", filter.Status, filter.ID, filter.SortPriority)
	var orders []types.ProductionOrder

	if services.GetCache(cacheKey, &orders) {
		return &orders, nil
	}

	query := db.Preload("Material").Preload("Machine")

	if filter.ID != 0 {
		query = query.Preload("Items").Where("id = ?", filter.ID)
	} else {
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
	}

	// Ordenamiento
	if filter.SortPriority {
		query = query.Order("priority ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	services.SetCache(cacheKey, &orders)
	return &orders, nil
}
