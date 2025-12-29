package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllOrdersUseCase(filter types.OrderFilter) (*[]types.ProductionOrder, error) {
	db := storage.DatabaseInstance{}.Instance()
	var orders []types.ProductionOrder

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
	return &orders, nil
}
