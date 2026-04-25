package services

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *StockService) UpdateItem(id uint, stock *types.ProductUpdateRequest, companyID uint) error {
	exists, err := s.repo.GetByID(&id, &stock.SKU, companyID)
	if err != nil {
		return fmt.Errorf("error getting stock: %w", err)
	}
	if exists == nil {
		return fmt.Errorf("the item does not exist: %d", id)
	}
	if stock.Name != "" {
		exists.Name = stock.Name
	}
	if stock.Description != "" {
		exists.Description = stock.Description
	}
	if stock.Quantity != 0 {
		exists.Quantity = stock.Quantity
	}
	if stock.Price >= 0 {
		exists.Price = stock.Price
	}
	if stock.MinQty >= 0 {
		exists.MinQty = stock.MinQty
	}
	return s.repo.Update(id, exists, companyID)
}
