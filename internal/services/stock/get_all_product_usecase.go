package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllProductsUseCase() (*[]types.StockItem, error) {
	db := storage.DatabaseInstance{}.Instance()
	cacheKey := "stock:list:all"
	var items []types.StockItem
	if services.GetCache(cacheKey, &items) {
		return &items, nil
	}

	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}
	services.SetCache(cacheKey, &items)
	return &items, nil
}
