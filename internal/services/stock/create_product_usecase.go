package services

import (
	"fmt"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/database/stock_utils"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func CreateProductUseCase(item types.ProductCreateRequest) error {
	db := storage.DatabaseInstance{}.Instance()

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

	product := types.StockItem{
		SKU:         item.SKU,
		Name:        item.Name,
		Description: item.Description,
		Quantity:    item.Quantity,
		Price:       item.Price,
	}
	if err := db.Create(&product).Error; err != nil {
		return fmt.Errorf("The Product already exists")
	}
	return nil
}
