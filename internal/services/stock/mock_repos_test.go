package services

import (
	"errors"

	"github.com/francotraversa/Sliceflow/internal/types"
)

// --- Mock: StockRepository ---

type mockStockRepo struct {
	createFn       func(stock *types.StockItem) error
	getByIDFn      func(id *uint, sku *string, companyID uint) (*types.StockItem, error)
	getAllFn       func(companyID uint) (*[]types.StockItem, error)
	updateFn       func(id uint, stock *types.StockItem, companyID uint) error
	deleteFn       func(id uint, companyID uint) error
	getDashboardFn func(companyID uint) (*types.DashboardResponse, error)
}

func (m *mockStockRepo) Create(s *types.StockItem) error {
	if m.createFn != nil {
		return m.createFn(s)
	}
	return nil
}

func (m *mockStockRepo) GetByID(id *uint, sku *string, companyID uint) (*types.StockItem, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(id, sku, companyID)
	}
	return nil, nil
}

func (m *mockStockRepo) GetAll(companyID uint) (*[]types.StockItem, error) {
	if m.getAllFn != nil {
		return m.getAllFn(companyID)
	}
	return &[]types.StockItem{}, nil
}

func (m *mockStockRepo) Update(id uint, s *types.StockItem, companyID uint) error {
	if m.updateFn != nil {
		return m.updateFn(id, s, companyID)
	}
	return nil
}

func (m *mockStockRepo) Delete(id uint, companyID uint) error {
	if m.deleteFn != nil {
		return m.deleteFn(id, companyID)
	}
	return nil
}

func (m *mockStockRepo) GetDashboardStats(companyID uint) (*types.DashboardResponse, error) {
	if m.getDashboardFn != nil {
		return m.getDashboardFn(companyID)
	}
	return &types.DashboardResponse{}, nil
}

// --- Mock: StockAuditRepository ---

type mockAuditRepo struct {
	executeMovementFn func(req types.CreateMovementRequest, companyID uint) error
	getMovementByIDFn func(id uint, companyID uint) (*types.StockMovement, error)
	getAllMovementsFn func(filter types.HistoryFilter, companyID uint) ([]types.StockMovement, error)
}

func (m *mockAuditRepo) ExecuteMovementTransaction(req types.CreateMovementRequest, companyID uint) error {
	if m.executeMovementFn != nil {
		return m.executeMovementFn(req, companyID)
	}
	return nil
}

func (m *mockAuditRepo) GetMovementByID(id uint, companyID uint) (*types.StockMovement, error) {
	if m.getMovementByIDFn != nil {
		return m.getMovementByIDFn(id, companyID)
	}
	return nil, nil
}

func (m *mockAuditRepo) GetAllMovements(filter types.HistoryFilter, companyID uint) ([]types.StockMovement, error) {
	if m.getAllMovementsFn != nil {
		return m.getAllMovementsFn(filter, companyID)
	}
	return []types.StockMovement{}, nil
}

// --- Helpers ---

var errRepoDown = errors.New("repository unavailable")
