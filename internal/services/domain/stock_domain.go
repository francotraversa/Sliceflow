package domain

import "github.com/francotraversa/Sliceflow/internal/types"

type StockService interface {
	CreateItem(stock *types.ProductCreateRequest, companyID uint) error
	UpdateItem(id uint, stock *types.ProductUpdateRequest, companyID uint) error
	DeleteItem(id uint, companyID uint) error
	GetItemByID(id uint, companyID uint) (*types.StockItem, error)
	GetAllItems(companyID uint) (*[]types.StockItem, error)
	GetDashboard(companyID uint) (*types.DashboardResponse, error)
}

type StockRepository interface {
	Create(stock *types.StockItem) error
	Update(id uint, stock *types.StockItem, companyID uint) error
	Delete(id uint, companyID uint) error
	GetByID(id *uint, sku *string, companyID uint) (*types.StockItem, error)
	GetAll(companyID uint) (*[]types.StockItem, error)
	GetDashboardStats(companyID uint) (*types.DashboardResponse, error)
}
