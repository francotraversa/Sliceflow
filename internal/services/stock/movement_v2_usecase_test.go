package services

import (
	"testing"

	"github.com/francotraversa/Sliceflow/internal/types"
)

// --- CreateMovement ---

func TestCreateMovement_MissingSKU(t *testing.T) {
	svc := NewStockAuditService(&mockAuditRepo{})
	err := svc.CreateMovement(types.CreateMovementRequest{Type: "IN", Quantity: 5, UserID: 1}, testCompany)
	if err == nil {
		t.Fatal("expected error for missing SKU")
	}
}

func TestCreateMovement_MissingType(t *testing.T) {
	svc := NewStockAuditService(&mockAuditRepo{})
	err := svc.CreateMovement(types.CreateMovementRequest{SKU: "X", Quantity: 5, UserID: 1}, testCompany)
	if err == nil {
		t.Fatal("expected error for missing Type")
	}
}

func TestCreateMovement_ZeroQuantity(t *testing.T) {
	svc := NewStockAuditService(&mockAuditRepo{})
	err := svc.CreateMovement(types.CreateMovementRequest{SKU: "X", Type: "IN", UserID: 1}, testCompany)
	if err == nil {
		t.Fatal("expected error for zero Quantity")
	}
}

func TestCreateMovement_MissingUserID(t *testing.T) {
	svc := NewStockAuditService(&mockAuditRepo{})
	err := svc.CreateMovement(types.CreateMovementRequest{SKU: "X", Type: "IN", Quantity: 5}, testCompany)
	if err == nil {
		t.Fatal("expected error for missing UserID")
	}
}

func TestCreateMovement_RepoError(t *testing.T) {
	svc := NewStockAuditService(&mockAuditRepo{
		executeMovementFn: func(req types.CreateMovementRequest, companyID uint) error {
			return errRepoDown
		},
	})
	err := svc.CreateMovement(types.CreateMovementRequest{
		SKU:      "ABC",
		Type:     "IN",
		Quantity: 10,
		UserID:   1,
	}, testCompany)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCreateMovement_Success(t *testing.T) {
	called := false
	svc := NewStockAuditService(&mockAuditRepo{
		executeMovementFn: func(req types.CreateMovementRequest, companyID uint) error {
			called = true
			if req.SKU != "ABC" {
				return errRepoDown
			}
			return nil
		},
	})
	err := svc.CreateMovement(types.CreateMovementRequest{
		SKU:        "ABC",
		Type:       "IN",
		Quantity:   10,
		UserID:     1,
		LocationID: 1,
		Reason:     "Compra",
		Description: "Test",
	}, testCompany)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("repo.ExecuteMovementTransaction was never called")
	}
}

// --- GetMovementByID ---

func TestGetMovementByID_NotFound(t *testing.T) {
	svc := NewStockAuditService(&mockAuditRepo{
		getMovementByIDFn: func(id uint, companyID uint) (*types.StockMovement, error) {
			return nil, nil
		},
	})
	mov, err := svc.GetMovementByID(99, testCompany)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mov != nil {
		t.Error("expected nil movement")
	}
}

func TestGetMovementByID_Success(t *testing.T) {
	svc := NewStockAuditService(&mockAuditRepo{
		getMovementByIDFn: func(id uint, companyID uint) (*types.StockMovement, error) {
			return &types.StockMovement{ID: id, Type: "IN"}, nil
		},
	})
	mov, err := svc.GetMovementByID(5, testCompany)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mov.ID != 5 {
		t.Errorf("wrong ID: %d", mov.ID)
	}
}

// --- GetAllMovements ---

func TestGetAllMovements_FilterBySKU(t *testing.T) {
	allMovements := []types.StockMovement{
		{ID: 1, StockSKU: "AAA", Type: "IN"},
		{ID: 2, StockSKU: "BBB", Type: "OUT"},
		{ID: 3, StockSKU: "AAA", Type: "OUT"},
	}

	svc := NewStockAuditService(&mockAuditRepo{
		getAllMovementsFn: func(filter types.HistoryFilter, companyID uint) ([]types.StockMovement, error) {
			// simula filtrado del repo
			var result []types.StockMovement
			for _, m := range allMovements {
				if filter.SKU == "" || m.StockSKU == filter.SKU {
					result = append(result, m)
				}
			}
			return result, nil
		},
	})

	results, err := svc.GetAllMovements(types.HistoryFilter{SKU: "AAA"}, testCompany)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 movements for SKU AAA, got %d", len(results))
	}
	for _, m := range results {
		if m.StockSKU != "AAA" {
			t.Errorf("filter leaked SKU %s", m.StockSKU)
		}
	}
}

func TestGetAllMovements_Empty(t *testing.T) {
	svc := NewStockAuditService(&mockAuditRepo{
		getAllMovementsFn: func(filter types.HistoryFilter, companyID uint) ([]types.StockMovement, error) {
			return []types.StockMovement{}, nil
		},
	})
	results, err := svc.GetAllMovements(types.HistoryFilter{}, testCompany)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty slice, got %d items", len(results))
	}
}
