package services

import (
	"fmt"

	servicesWeb "github.com/francotraversa/Sliceflow/internal/services/common"
)

func (s *OrderService) DeleteOrder(id uint, companyID uint) error {
	if err := s.repo.DeleteOrder(id, companyID); err != nil {
		return fmt.Errorf("could not delete order: %w", err)
	}
	servicesWeb.PublishEvent("dashboard_updates", `{"type": "ORDER_DELETED", "message": "ORDER DELETED"}`)
	return nil
}
