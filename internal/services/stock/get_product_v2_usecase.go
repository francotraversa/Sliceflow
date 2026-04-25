package services

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *StockService) GetItemByID(id uint, companyID uint) (*types.StockItem, error) {
	exists, err := s.repo.GetByID(&id, nil, companyID)
	if err != nil {
		return nil, fmt.Errorf("error getting stock: %w", err)
	}
	if exists == nil {
		return nil, fmt.Errorf("the item does not exist: %d", id)
	}
	return exists, nil
}

func (s *StockService) GetItemBySKU(sku string, companyID uint) (*types.StockItem, error) {
	exists, err := s.repo.GetByID(nil, &sku, companyID)
	if err != nil {
		return nil, fmt.Errorf("error getting stock: %w", err)
	}
	if exists == nil {
		return nil, fmt.Errorf("the item does not exist: %s", sku)
	}
	return exists, nil
}

func (s *StockService) GetAllItems(companyID uint) (*[]types.StockItem, error) {
	stocks, err := s.repo.GetAll(companyID)
	if err != nil {
		return nil, fmt.Errorf("error getting stocks: %w", err)
	}
	return stocks, nil
}
