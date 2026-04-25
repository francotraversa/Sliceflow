package services

import (
	"testing"

	"github.com/francotraversa/Sliceflow/internal/types"
)

const testCompany uint = 1

// --- CreateItem ---

func TestCreateItem_MissingSKU(t *testing.T) {
	svc := NewStockService(&mockStockRepo{})
	err := svc.CreateItem(&types.ProductCreateRequest{Name: "Test"}, testCompany)
	if err == nil {
		t.Fatal("expected error for missing SKU")
	}
}

func TestCreateItem_MissingName(t *testing.T) {
	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(id *uint, sku *string, companyID uint) (*types.StockItem, error) {
			return nil, nil // SKU no existe aún
		},
	})
	err := svc.CreateItem(&types.ProductCreateRequest{SKU: "ABC-001"}, testCompany)
	if err == nil {
		t.Fatal("expected error for missing Name")
	}
}

func TestCreateItem_DuplicateSKU(t *testing.T) {
	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(id *uint, sku *string, companyID uint) (*types.StockItem, error) {
			return &types.StockItem{SKU: *sku}, nil // simula que ya existe
		},
	})
	err := svc.CreateItem(&types.ProductCreateRequest{SKU: "ABC-001", Name: "Dup"}, testCompany)
	if err == nil {
		t.Fatal("expected error for duplicate SKU")
	}
}

func TestCreateItem_NegativePrice(t *testing.T) {
	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(_ *uint, _ *string, _ uint) (*types.StockItem, error) { return nil, nil },
	})
	err := svc.CreateItem(&types.ProductCreateRequest{SKU: "X", Name: "Y", Price: -1}, testCompany)
	if err == nil {
		t.Fatal("expected error for negative price")
	}
}

func TestCreateItem_RepoError(t *testing.T) {
	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(_ *uint, _ *string, _ uint) (*types.StockItem, error) {
			return nil, errRepoDown
		},
	})
	err := svc.CreateItem(&types.ProductCreateRequest{SKU: "X", Name: "Y"}, testCompany)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCreateItem_Success(t *testing.T) {
	var createdItem *types.StockItem
	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(_ *uint, _ *string, _ uint) (*types.StockItem, error) { return nil, nil },
		createFn: func(s *types.StockItem) error {
			createdItem = s
			return nil
		},
	})
	err := svc.CreateItem(&types.ProductCreateRequest{
		SKU:      "NEW-001",
		Name:     "Tornillo",
		Quantity: 100,
		Price:    10,
	}, testCompany)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if createdItem.SKU != "NEW-001" {
		t.Errorf("wrong SKU saved: %s", createdItem.SKU)
	}
	if createdItem.IdCompany != testCompany {
		t.Errorf("companyID not set correctly")
	}
}

// --- DeleteItem ---

func TestDeleteItem_NotFound(t *testing.T) {
	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(_ *uint, _ *string, _ uint) (*types.StockItem, error) { return nil, nil },
	})
	err := svc.DeleteItem(99, testCompany)
	if err == nil {
		t.Fatal("expected error when item not found")
	}
}

func TestDeleteItem_Success(t *testing.T) {
	deleted := false
	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(_ *uint, _ *string, _ uint) (*types.StockItem, error) {
			return &types.StockItem{Id: 1}, nil
		},
		deleteFn: func(id uint, companyID uint) error {
			deleted = true
			return nil
		},
	})
	err := svc.DeleteItem(1, testCompany)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("repo.Delete was never called")
	}
}

// --- UpdateItem ---

func TestUpdateItem_NotFound(t *testing.T) {
	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(_ *uint, _ *string, _ uint) (*types.StockItem, error) { return nil, nil },
	})
	err := svc.UpdateItem(1, &types.ProductUpdateRequest{SKU: "X"}, testCompany)
	if err == nil {
		t.Fatal("expected error when item not found")
	}
}

func TestUpdateItem_Success(t *testing.T) {
	original := &types.StockItem{Id: 1, SKU: "ABC", Name: "Old", Quantity: 10}
	var saved *types.StockItem

	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(_ *uint, _ *string, _ uint) (*types.StockItem, error) { return original, nil },
		updateFn: func(id uint, s *types.StockItem, companyID uint) error {
			saved = s
			return nil
		},
	})

	err := svc.UpdateItem(1, &types.ProductUpdateRequest{SKU: "ABC", Name: "New Name"}, testCompany)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved.Name != "New Name" {
		t.Errorf("name not updated: %s", saved.Name)
	}
}

// --- GetItemByID ---

func TestGetItemByID_NotFound(t *testing.T) {
	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(_ *uint, _ *string, _ uint) (*types.StockItem, error) { return nil, nil },
	})
	_, err := svc.GetItemByID(1, testCompany)
	if err == nil {
		t.Fatal("expected error when item not found")
	}
}

func TestGetItemByID_Success(t *testing.T) {
	svc := NewStockService(&mockStockRepo{
		getByIDFn: func(_ *uint, _ *string, _ uint) (*types.StockItem, error) {
			return &types.StockItem{Id: 1, Name: "Tornillo"}, nil
		},
	})
	item, err := svc.GetItemByID(1, testCompany)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Name != "Tornillo" {
		t.Errorf("wrong name: %s", item.Name)
	}
}

// --- GetDashboard ---

func TestGetDashboard_ReturnsData(t *testing.T) {
	expected := &types.DashboardResponse{TotalItems: 42, TotalValue: 9999.99}
	svc := NewStockService(&mockStockRepo{
		getDashboardFn: func(companyID uint) (*types.DashboardResponse, error) {
			return expected, nil
		},
	})
	res, err := svc.GetDashboard(testCompany)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TotalItems != 42 {
		t.Errorf("wrong TotalItems: %d", res.TotalItems)
	}
}
