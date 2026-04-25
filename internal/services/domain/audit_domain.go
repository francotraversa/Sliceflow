package domain

import "github.com/francotraversa/Sliceflow/internal/types"

type StockAuditService interface {
	CreateMovement(req types.CreateMovementRequest, companyID uint) error
	GetMovementByID(id uint, companyID uint) (*types.StockMovement, error)
	GetAllMovements(filter types.HistoryFilter, companyID uint) ([]types.StockMovement, error)
}

type StockAuditRepository interface {
	ExecuteMovementTransaction(req types.CreateMovementRequest, companyID uint) error
	GetMovementByID(id uint, companyID uint) (*types.StockMovement, error)
	GetAllMovements(filter types.HistoryFilter, companyID uint) ([]types.StockMovement, error)
}
