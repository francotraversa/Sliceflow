package services

import (
	"fmt"
	"strings"

	"github.com/francotraversa/Sliceflow/internal/infra/database/stock_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func CreateProductUseCase(item types.ProductCreateRequest) error {
	item.SKU = strings.ToUpper(strings.TrimSpace(item.SKU))
	item.Name = strings.TrimSpace(item.Name)

	if item.SKU == "" || item.Name == "" {
		return fmt.Errorf("SKU and Name must be filled")
	}
	exist, err := stock_utils.CheckProductExistsBySKU(item.SKU)

	if err != nil {
		return err
	}

	if exist != nil {
		return fmt.Errorf("the product with SKU %s already exists", item.SKU)
	}

	if item.Price < 0 {
		return fmt.Errorf("The price must be positive")
	}

	if item.MinQty < 0 {
		return fmt.Errorf("The minimum quantity must be positive")
	}

	if item.Quantity < 0 {
		return fmt.Errorf("The quantity must be positive")
	}

	if item.Description == "" {
		item.Description = "No description"
	}

	if item.MinQty == 0 {
		item.MinQty = 5 // Valor por defecto si no se proporciona
	}

	product := types.StockItem{
		SKU:         item.SKU,
		Name:        item.Name,
		Description: item.Description,
		Quantity:    item.Quantity,
		Price:       item.Price,
		MinQty:      item.MinQty,
	}
	if err := db_utils.Create(&product); err != nil {
		return fmt.Errorf("Error Create Product")
	}
	services.InvalidateCache("stock:list:all")
	services.PublishEvent("dashboard_updates", `{"type": "PRODUCT_CREATED", "message": "PRODUCT CREATED"}`)
	return nil
}
