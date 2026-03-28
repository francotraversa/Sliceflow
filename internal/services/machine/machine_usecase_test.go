package services

import (
	"testing"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testCompanyIDMachine uint = 1

// setupMachineTest creates a clean in-memory DB for each test
func setupMachineTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	// Migrate Machine table
	err = db.AutoMigrate(&types.Machine{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Inject test DB into global instance
	storage.OverrideDatabaseInstance(db)
	return db
}

func TestMachineCRUD(t *testing.T) {

	// --- 1. CREATE TESTS ---
	t.Run("Create Machine Successfully", func(t *testing.T) {
		db := setupMachineTest(t)

		dto := types.CreateMachineDTO{
			Name: "Prusa MK3S - Sector A",
			Type: "FDM",
		}

		err := CreateMachineUseCase(dto, testCompanyIDMachine)
		if err != nil {
			t.Fatalf("Unexpected error creating machine: %v", err)
		}

		// Verify it was saved and default status is 'idle'
		var machine types.Machine
		db.First(&machine, "name = ?", dto.Name)

		if machine.ID == 0 {
			t.Error("Machine was not saved to the DB")
		}
		if machine.Status != "idle" {
			t.Errorf("Initial status should be 'idle', got: %s", machine.Status)
		}
	})

	t.Run("Create Duplicate Machine (Should Fail)", func(t *testing.T) {
		setupMachineTest(t)

		dto := types.CreateMachineDTO{Name: "Ender 3 V2", Type: "FDM"}

		// First insert
		CreateMachineUseCase(dto, testCompanyIDMachine)

		// Second insert (Same name)
		err := CreateMachineUseCase(dto, testCompanyIDMachine)
		if err == nil {
			t.Error("Should have failed due to duplicate name, but passed")
		}
	})

	// --- 2. READ TESTS ---
	t.Run("List Machines", func(t *testing.T) {
		setupMachineTest(t)

		CreateMachineUseCase(types.CreateMachineDTO{Name: "M1", Type: "FDM"}, testCompanyIDMachine)
		CreateMachineUseCase(types.CreateMachineDTO{Name: "M2", Type: "SLS"}, testCompanyIDMachine)

		list, err := GetAllMachinesUseCase(types.MachineFilter{}, testCompanyIDMachine)
		if err != nil {
			t.Fatalf("Failed to list: %v", err)
		}

		if len(*list) != 2 {
			t.Errorf("Expected 2 machines, got %d", len(*list))
		}
	})

	// --- 3. UPDATE TESTS ---
	t.Run("Update Machine", func(t *testing.T) {
		db := setupMachineTest(t)

		// Create original
		CreateMachineUseCase(types.CreateMachineDTO{Name: "Original"}, testCompanyIDMachine)

		// Update
		id := 1
		editada := "Edited"
		sla := "SLA"
		maintenance := "maintenance"
		updateDTO := types.UpdateMachineDTO{
			Name:   &editada,
			Type:   &sla,
			Status: &maintenance, // Test manual status change
		}

		err := UpdateMachineUseCase(id, updateDTO, testCompanyIDMachine)
		if err != nil {
			t.Fatalf("Failed to update: %v", err)
		}

		// Verify
		var m types.Machine
		db.First(&m, id)

		if m.Name != "Edited" || m.Status != "maintenance" {
			t.Errorf("Incorrect data. Name: %s, Status: %s", m.Name, m.Status)
		}
	})

	// --- 4. DELETE (SOFT DELETE) TESTS ---
	t.Run("Delete Machine (Soft Delete)", func(t *testing.T) {
		db := setupMachineTest(t)

		// Create
		CreateMachineUseCase(types.CreateMachineDTO{Name: "To Delete", Type: "FDM"}, testCompanyIDMachine)

		// Delete ID 1
		err := DeleteMachineUseCase(1, testCompanyIDMachine)
		if err != nil {
			t.Fatalf("Failed to delete: %v", err)
		}

		// A. Verify GetAll does NOT return it
		list, _ := GetAllMachinesUseCase(types.MachineFilter{}, testCompanyIDMachine)

		if len(*list) != 0 {
			t.Errorf("List should be empty, got %d", len(*list))
		}

		// B. Verify it still exists in DB with deleted_at (Unscoped)
		var count int64
		db.Unscoped().Model(&types.Machine{}).Where("id = ?", 1).Count(&count)
		if count != 1 {
			t.Error("Physical record disappeared (should have been soft deleted)")
		}
	})
}
