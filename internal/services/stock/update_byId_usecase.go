package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/database/stock_utils"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func UpdateByIdProductUseCase(sku string, req types.ProductUpdateRequest) (*types.StockItem, error) {
	db := storage.DatabaseInstance{}.Instance()
	itemexist, err := stock_utils.CheckProductExistsBySKU(sku)

	if itemexist == nil {
		return nil, fmt.Errorf("The product doesn't exist")
	}
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		itemexist.Name = req.Name
	}
	if req.Description != "" {
		itemexist.Description = req.Description
	}
	if req.Quantity != 0 {
		itemexist.Quantity = req.Quantity
	}

	if err := db.Save(&itemexist).Error; err != nil {
		return nil, err
	}
	return itemexist, nil
}
