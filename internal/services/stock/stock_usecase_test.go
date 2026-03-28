package services

import (
	"testing"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testCompanyIDStock uint = 1

// setupStockTest sets up an in-memory DB for product catalog tests
func setupStockTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	// Migrate the Stock model
	db.AutoMigrate(&types.StockItem{})

	// Inject the test DB into the global instance
	storage.OverrideDatabaseInstance(db)
	return db
}
func TestStockService(t *testing.T) {
	db := setupStockTest(t)

	t.Run("Create Product Successfully", func(t *testing.T) {
		req := types.ProductCreateRequest{
			SKU:         "SKU-779-1",
			Name:        "Phillips Screw",
			Description: "1/2 inch",
			Quantity:    100,
			Price:       90,
		}

		err := CreateProductUseCase(req, testCompanyIDStock)
		if err != nil {
			t.Fatalf("Should not have errored when creating product: %v", err)
		}

		var found types.StockItem
		db.Where("sku = ?", "SKU-779-1").First(&found)
		if found.Name != "Phillips Screw" {
			t.Errorf("Expected 'Phillips Screw', got %s", found.Name)
		}
	})

	t.Run("Prevent Duplicate SKU", func(t *testing.T) {
		req := types.ProductCreateRequest{
			SKU:  "SKU-779-1", // SKU already exists from previous test
			Name: "Duplicate Product",
		}

		err := CreateProductUseCase(req, testCompanyIDStock)
		if err == nil {
			t.Error("Expected error for duplicate SKU")
		}
	})

	t.Run("Soft Delete Product", func(t *testing.T) {
		// Create a new product to delete
		item := types.StockItem{SKU: "DEL-99", Name: "DeleteMe", Status: "active", IdCompany: testCompanyIDStock}
		db.Create(&item)

		// Call the UseCase
		err := DeleteByIdUseCase(item.SKU, testCompanyIDStock)
		if err != nil {
			t.Fatalf("Failed to delete: %v", err)
		}

		// 1. Verify it's NOT found in normal queries
		var found types.StockItem
		result := db.Where("sku = ?", item.SKU).First(&found)

		if result.Error == nil {
			t.Error("Product is still visible in normal queries")
		}

		// 2. Verify it STILL exists physically in the DB (Unscoped)
		var raw types.StockItem
		db.Unscoped().Where("sku = ?", item.SKU).First(&raw)

		if !raw.DeletedAt.Valid {
			t.Error("DeletedAt was not set by GORM")
		}
	})

	t.Run("Update Product - Success and Failure", func(t *testing.T) {
		// A. Happy Path
		item := types.StockItem{SKU: "EDIT-1", Name: "Old Name", Quantity: 50, IdCompany: testCompanyIDStock}
		db.Create(&item)

		updateReq := types.ProductUpdateRequest{Name: "New Name"}

		// Call UseCase with SKU string
		_, err := UpdateByIdProductUseCase(item.SKU, updateReq, testCompanyIDStock)
		if err != nil {
			t.Errorf("Unexpected error during update: %v", err)
		}

		// Validate the change
		var found types.StockItem
		db.Where("sku = ?", item.SKU).First(&found)

		if found.Name != "New Name" {
			t.Errorf("Name was not updated. Current value: %s", found.Name)
		}

		// B. Failure Case (SKU does not exist)
		_, err = UpdateByIdProductUseCase("FAKE-SKU-999", updateReq, testCompanyIDStock)
		if err == nil {
			t.Error("Expected error when updating non-existent SKU, but it didn't fail")
		}
	})
}
