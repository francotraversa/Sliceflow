package services

import (
	"errors"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetProductByNameUseCase(name string) (*[]types.StockItem, error) {
	db := storage.DatabaseInstance{}.Instance()
	var items []types.StockItem
	if err := db.Where("name = ?", name).Find(&items).Error; err != nil {
		return nil, errors.New("product not found")
	}
	return &items, nil
}
