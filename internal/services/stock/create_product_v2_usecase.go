package services

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *StockService) CreateItem(stock *types.ProductCreateRequest, companyID uint) error {
	if stock.SKU == "" {
		return fmt.Errorf("sku is required")
	}
	exists, err := s.repo.GetByID(nil, &stock.SKU, companyID)
	if err != nil {
		return fmt.Errorf("error getting stock: %w", err)
	}
	if exists != nil {
		return fmt.Errorf("the item already exists: %s", stock.SKU)
	}
	if stock.Name == "" {
		return fmt.Errorf("name is required")
	}
	if stock.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}
	if stock.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	if stock.MinQty < 0 {
		return fmt.Errorf("min_qty cannot be negative")
	}

	newItem := types.StockItem{
		SKU:         stock.SKU,
		Name:        stock.Name,
		Description: stock.Description,
		Quantity:    stock.Quantity,
		Price:       stock.Price,
		MinQty:      stock.MinQty,
		IdCompany:   companyID,
	}

	return s.repo.Create(&newItem)
}
