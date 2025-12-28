package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetStockHistoryUseCase(filter types.HistoryFilter) ([]types.StockMovement, error) {
	db := storage.DatabaseInstance{}.Instance()
	var movements []types.StockMovement

	query := db.Model(&types.StockMovement{})

	if filter.SKU != "" {
		query = query.Where("stock_sku = ?", filter.SKU)
	}

	if filter.StartDate != "" {
		query = query.Where("created_at >= ?", filter.StartDate+" 00:00:00")
	}

	if filter.EndDate != "" {
		query = query.Where("created_at <= ?", filter.EndDate+" 23:59:59")
	}

	result := query.Order("created_at desc").Find(&movements)

	if result.Error != nil {
		return nil, result.Error
	}

	return movements, nil
}
