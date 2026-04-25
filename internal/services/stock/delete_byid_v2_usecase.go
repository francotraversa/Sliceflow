package services

import (
	"fmt"
)

func (s *StockService) DeleteItem(id uint, companyID uint) error {
	exists, err := s.repo.GetByID(&id, nil, companyID)
	if err != nil {
		return fmt.Errorf("error getting stock: %w", err)
	}
	if exists == nil {
		return fmt.Errorf("the item does not exist: %d", id)
	}
	return s.repo.Delete(id, companyID)
}
