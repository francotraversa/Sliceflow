package services

import "github.com/francotraversa/Sliceflow/internal/types"

func (s *OrderService) DashboardOrders(userRole string, companyID uint) (*types.ProductionDashboardResponse, error) {
	response, err := s.repo.DashboardOrders(companyID)
	if err != nil {
		return nil, err
	}

	isAdmin := (userRole == "admin")

	if !isAdmin {
		response.TotalRevenueFDM = 0

		censoredOrders := make([]types.ProductionOrder, len(response.Orders))
		copy(censoredOrders, response.Orders)

		for i := range censoredOrders {
			censoredOrders[i].TotalPrice = nil
		}
		response.Orders = censoredOrders
	}

	return response, nil
}
