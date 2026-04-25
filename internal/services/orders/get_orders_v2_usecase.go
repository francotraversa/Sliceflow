package services

import (
	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *OrderService) GetOrdersByStatus(filter types.OrderFilter, companyID uint) (*[]types.ProductionOrder, error) {
	orders, err := s.repo.GetOrdersByFilter(filter, companyID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
