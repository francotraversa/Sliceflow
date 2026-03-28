package services

import (
	"testing"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- SETUP & HELPERS ---

const testCompanyIDMovement uint = 1

// setupMovementTest starts a clean in-memory DB for each test
func setupMovementTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	// Migrate StockItem and StockMovement
	db.AutoMigrate(&types.StockItem{}, &types.StockMovement{})

	storage.OverrideDatabaseInstance(db)
	return db
}

// parseDate helps create specific dates for filter tests
func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}

// --- TEST 1: WRITE OPERATIONS (MOVEMENTS) ---

func TestStockMovements(t *testing.T) {
	db := setupMovementTest(t)

	// Prepare a base product for the tests
	baseSku := "TEST-MOV-1"
	initialItem := types.StockItem{
		SKU:       baseSku,
		Name:      "Test Product",
		Quantity:  10, // Start with 10 units
		Price:     100,
		IdCompany: testCompanyIDMovement,
	}
	db.Create(&initialItem)

	t.Run("Inbound Movement (IN) - Adds Stock", func(t *testing.T) {
		req := types.CreateMovementRequest{
			SKU:      baseSku,
			Type:     "IN",
			Quantity: 5, // 10 + 5 = 15
			UserID:   1,
			Reason:   "Supplier Purchase",
		}

		err := AddStockMovementUseCase(req, testCompanyIDMovement)
		if err != nil {
			t.Fatalf("Unexpected error on IN movement: %v", err)
		}

		// 1. Verify the StockItem was updated
		var item types.StockItem
		db.First(&item, "sku = ?", baseSku)
		if item.Quantity != 15 {
			t.Errorf("Stock should have increased to 15, but is %d", item.Quantity)
		}

		// 2. Verify a history record was created (Audit)
		var mov types.StockMovement
		result := db.Where("stock_sku = ?", baseSku).Last(&mov)
		if result.Error != nil {
			t.Fatalf("StockMovement record was not created")
		}

		// Validate the internal math
		if mov.QtyBefore != 10 {
			t.Errorf("QtyBefore incorrect. Expected 10, got %d", mov.QtyBefore)
		}
		if mov.QtyDelta != 5 {
			t.Errorf("QtyDelta incorrect. Expected 5, got %d", mov.QtyDelta)
		}
		if mov.QtyAfter != 15 {
			t.Errorf("QtyAfter incorrect. Expected 15, got %d", mov.QtyAfter)
		}
	})

	t.Run("Outbound Movement (OUT) - Reduces Stock", func(t *testing.T) {
		// Starting from 15 (result of previous test)
		req := types.CreateMovementRequest{
			SKU:      baseSku,
			Type:     "OUT",
			Quantity: 3, // 15 - 3 = 12
			UserID:   1,
			Reason:   "Customer Sale",
		}

		err := AddStockMovementUseCase(req, testCompanyIDMovement)
		if err != nil {
			t.Fatalf("Unexpected error on OUT movement: %v", err)
		}

		// 1. Verify Stock
		var item types.StockItem
		db.First(&item, "sku = ?", baseSku)
		if item.Quantity != 12 {
			t.Errorf("Stock should have decreased to 12, but is %d", item.Quantity)
		}

		// 2. Verify Audit
		var mov types.StockMovement
		db.Where("stock_sku = ?", baseSku).Last(&mov)

		if mov.QtyBefore != 15 {
			t.Errorf("QtyBefore wrong: %d", mov.QtyBefore)
		}
		if mov.QtyDelta != -3 {
			t.Errorf("QtyDelta wrong (should be negative): %d", mov.QtyDelta)
		}
		if mov.QtyAfter != 12 {
			t.Errorf("QtyAfter wrong: %d", mov.QtyAfter)
		}
	})

	t.Run("Validate Insufficient Stock (Error)", func(t *testing.T) {
		req := types.CreateMovementRequest{
			SKU:      baseSku,
			Type:     "OUT",
			Quantity: 100, // Impossible (only 12 available)
			UserID:   1,
		}

		err := AddStockMovementUseCase(req, testCompanyIDMovement)
		if err == nil {
			t.Error("Expected insufficient stock error, but operation succeeded")
		}

		// Verify stock was NOT modified
		var item types.StockItem
		db.First(&item, "sku = ?", baseSku)
		if item.Quantity != 12 {
			t.Errorf("Stock changed to %d despite the error", item.Quantity)
		}
	})

	t.Run("Fail on Non-Existent Product", func(t *testing.T) {
		req := types.CreateMovementRequest{
			SKU:      "GHOST-SKU",
			Type:     "IN",
			Quantity: 10,
		}

		err := AddStockMovementUseCase(req, testCompanyIDMovement)
		if err == nil {
			t.Error("Expected 'product not found' error")
		}
	})
}

// --- TEST 2: READ OPERATIONS (HISTORY & FILTERS) ---

func TestGetStockHistory(t *testing.T) {
	db := setupMovementTest(t) // Fresh clean DB for this test

	// --- SEEDING (Test data) ---
	skuA := "PROD-A"
	skuB := "PROD-B"

	// Insert movements manually with specific dates
	movements := []types.StockMovement{
		// January - Product A
		{StockSKU: skuA, Type: "IN", QtyDelta: 10, CreatedAt: parseDate("2025-01-01"), IdCompany: testCompanyIDMovement},
		{StockSKU: skuA, Type: "OUT", QtyDelta: -2, CreatedAt: parseDate("2025-01-15"), IdCompany: testCompanyIDMovement},

		// February - Product A
		{StockSKU: skuA, Type: "IN", QtyDelta: 5, CreatedAt: parseDate("2025-02-10"), IdCompany: testCompanyIDMovement},

		// January - Product B
		{StockSKU: skuB, Type: "IN", QtyDelta: 100, CreatedAt: parseDate("2025-01-05"), IdCompany: testCompanyIDMovement},
	}
	db.Create(&movements)

	// --- EXECUTION ---

	t.Run("Filter by SKU Only", func(t *testing.T) {
		// Request only Product A
		filter := types.HistoryFilter{SKU: skuA}

		results, err := GetStockHistoryUseCase(filter, testCompanyIDMovement)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		if len(*results) != 3 {
			t.Errorf("Expected 3 movements for SKU A, got %d", len(*results))
		}
		// Verify nothing from B is included
		for _, m := range *results {
			if m.StockSKU != skuA {
				t.Errorf("SKU filter failed: %s", m.StockSKU)
			}
		}
	})

	t.Run("Filter by Date Range Only (January)", func(t *testing.T) {
		// Request all movements from January (A and B)
		filter := types.HistoryFilter{StartDate: "2025-01-01", EndDate: "2025-01-31"}

		results, err := GetStockHistoryUseCase(filter, testCompanyIDMovement)
		if err != nil {
			t.Fatal(err)
		}

		// Expected: 2 from A + 1 from B = 3 Total. (February is excluded)
		if len(*results) != 3 {
			t.Errorf("Expected 3 movements in January, got %d", len(*results))
		}
	})

	t.Run("Combined Filter (SKU A + January)", func(t *testing.T) {
		filter := types.HistoryFilter{
			SKU:       skuA,
			StartDate: "2025-01-01",
			EndDate:   "2025-01-31",
		}

		results, err := GetStockHistoryUseCase(filter, testCompanyIDMovement)
		if err != nil {
			t.Fatal(err)
		}

		// Expected: Only the 2 from A in January
		if len(*results) != 2 {
			t.Errorf("Expected 2 combined movements, got %d", len(*results))
		}
	})
}
