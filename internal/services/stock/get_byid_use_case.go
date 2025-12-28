package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetByIdUseCase(sku uint) (*types.StockItem, error) {
	db := storage.DatabaseInstance{}.Instance()
	var item types.StockItem

	if err := db.First(&item, sku).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
