package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetStockHistoryUseCase(filter types.HistoryFilter) (*[]types.StockMovement, error) {
	db := storage.DatabaseInstance{}.Instance()
	cacheKey := fmt.Sprintf("historic:list:%s:%s", filter.SKU, filter.Type)

	var movements []types.StockMovement

	if services.GetCache(cacheKey, &movements) {
		return &movements, nil
	}

	query := db.Model(&types.StockMovement{})

	if filter.SKU != "" {
		query = query.Where("stock_sku = ?", filter.SKU)
	}

	if filter.StartDate != "" {
		query = query.Where("created_at >= ?", filter.StartDate+" 00:00:00")
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.EndDate != "" {
		query = query.Where("created_at <= ?", filter.EndDate+" 23:59:59")
	}

	result := query.Order("created_at desc").Find(&movements)

	if result.Error != nil {
		return nil, result.Error
	}
	services.SetCache(cacheKey, &movements)
	return &movements, nil
}
