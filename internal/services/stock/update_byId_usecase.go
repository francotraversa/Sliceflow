package services

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/infra/database/stock_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func UpdateByIdProductUseCase(sku string, req types.ProductUpdateRequest) (*types.StockItem, error) {
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
	if req.Price >= 0 {
		itemexist.Price = req.Price
	}

	if err := db_utils.Save(itemexist); err != nil {
		return nil, fmt.Errorf("Error update %s", itemexist.SKU)
	}
	services.InvalidateCache("stock:list:all")
	services.PublishEvent("dashboard_updates", `{"type": "PRODUCT_UPDATED", "message": "PRODUCT UPDATED"}`)

	return itemexist, nil
}
