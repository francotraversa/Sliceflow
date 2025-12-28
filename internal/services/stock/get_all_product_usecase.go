package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllProductsUseCase() (*[]types.StockItem, error) {
	db := storage.DatabaseInstance{}.Instance()
	var items []types.StockItem
	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}
	return &items, nil
}
