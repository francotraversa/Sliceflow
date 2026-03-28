package services

import (
	"testing"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testCompanyIDMaterial uint = 1

// setupMaterialTest creates a clean in-memory DB for each test
func setupMaterialTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	// Migrate only the Materials table
	err = db.AutoMigrate(&types.Material{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Inject test DB into global instance
	storage.OverrideDatabaseInstance(db)
	return db
}

func TestMaterialCRUD(t *testing.T) {
	// --- 1. CREATE TESTS ---
	t.Run("Create Material Successfully", func(t *testing.T) {
		db := setupMaterialTest(t)

		dto := types.CreateMaterialDTO{
			Name:        "PLA Black Grilon",
			Type:        "Filament",
			Description: "1kg spool",
		}

		err := CreateMaterialUseCase(dto, testCompanyIDMaterial)
		if err != nil {
			t.Fatalf("Unexpected error creating material: %v", err)
		}

		// Verify it was saved
		var count int64
		db.Model(&types.Material{}).Count(&count)
		if count != 1 {
			t.Errorf("Expected 1 material, got %d", count)
		}
	})

	t.Run("Create Duplicate Material (Should Fail)", func(t *testing.T) {
		setupMaterialTest(t)

		dto := types.CreateMaterialDTO{Name: "PLA Unique", Type: "Filament"}

		// First insert
		_ = CreateMaterialUseCase(dto, testCompanyIDMaterial)

		// Second insert (Same name)
		err := CreateMaterialUseCase(dto, testCompanyIDMaterial)
		if err == nil {
			t.Error("Should have failed due to duplicate name, but passed")
		}
	})

	// --- 2. READ TESTS (LIST) ---
	t.Run("List Materials", func(t *testing.T) {
		setupMaterialTest(t)

		// Insert 2 materials
		CreateMaterialUseCase(types.CreateMaterialDTO{Name: "Mat 1", Type: "A"}, testCompanyIDMaterial)
		CreateMaterialUseCase(types.CreateMaterialDTO{Name: "Mat 2", Type: "B"}, testCompanyIDMaterial)

		result, err := GetAllMaterialsUseCase(types.MaterialFilter{}, testCompanyIDMaterial)
		if err != nil {
			t.Fatalf("Failed to list: %v", err)
		}

		if len(*result) != 2 {
			t.Errorf("Expected 2 materials, got %d", len(*result))
		}
	})

	// --- 3. UPDATE TESTS ---
	t.Run("Update Material", func(t *testing.T) {
		db := setupMaterialTest(t)

		CreateMaterialUseCase(types.CreateMaterialDTO{Name: "Old Name", Type: "Old"}, testCompanyIDMaterial)

		id := 1
		updateDTO := types.UpdateMaterialDTO{
			Name:        "New Name",
			Type:        "New Type",
			Description: "Edited",
		}

		err := UpdateMaterialUseCase(id, updateDTO, testCompanyIDMaterial)
		if err != nil {
			t.Fatalf("Failed to update: %v", err)
		}

		var mat types.Material
		db.First(&mat, id)
		if mat.Name != "New Name" {
			t.Errorf("Name was not updated. Value: %s", mat.Name)
		}
	})

	// --- 4. DELETE (SOFT DELETE) TESTS ---
	t.Run("Delete Material (Soft Delete)", func(t *testing.T) {
		db := setupMaterialTest(t)

		// Create material ID 1
		CreateMaterialUseCase(types.CreateMaterialDTO{Name: "To Delete", Type: "X"}, testCompanyIDMaterial)

		// Delete
		err := DeleteMaterialUseCase(1, testCompanyIDMaterial)
		if err != nil {
			t.Fatalf("Failed to delete: %v", err)
		}

		// Verifications:

		// A. GetAll should NOT return it
		list, _ := GetAllMaterialsUseCase(types.MaterialFilter{}, testCompanyIDMaterial)

		if len(*list) != 0 {
			t.Errorf("List should be empty, got %d", len(*list))
		}

		// B. Physical record STILL exists in DB (Soft Delete)
		var count int64
		db.Unscoped().Model(&types.Material{}).Where("id = ?", 1).Count(&count)
		if count != 1 {
			t.Error("Physical record disappeared from DB (should have been soft deleted)")
		}

		// C. Verify deleted_at is not null
		var mat types.Material
		db.Unscoped().First(&mat, 1)
		if !mat.DeletedAt.Valid {
			t.Error("DeletedAt field is empty")
		}
	})
}
