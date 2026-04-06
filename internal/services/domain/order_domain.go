package domain

import "github.com/francotraversa/Sliceflow/internal/types"

type OrderRepository interface {
	Create(order *types.ProductionOrder) error
	GetByID(id uint, companyID uint) (bool, error)
	GetOrderWithItems(id uint, companyID uint) (*types.ProductionOrder, error)
	UpdateFullOrder(order *types.ProductionOrder, newItems []types.OrderItem) error
	GetOrdersByFilter(filter types.OrderFilter, companyID uint) (*[]types.ProductionOrder, error)
	DeleteOrder(id uint, companyID uint) error
	DashboardOrders(companyID uint) (*types.ProductionDashboardResponse, error)
}

type OrderUseCase interface {
	CreateOrder(order types.CreateOrderDTO, idCompany uint) error
	UpdateOrder(id uint, order types.UpdateOrderDTO, idCompany uint) error
	GetOrdersByStatus(filter types.OrderFilter, companyID uint) (*[]types.ProductionOrder, error)
	DeleteOrder(id uint, companyID uint) error
	DashboardOrders(userRole string, companyID uint) (*types.ProductionDashboardResponse, error)
}
