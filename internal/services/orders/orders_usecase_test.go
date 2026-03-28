package services

import (
	"testing"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testCompanyIDOrder uint = 1

// setupOrderTest: Prepares an in-memory DB and required tables
func setupOrderTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	// Migrate ALL tables that the use case touches, including OrderItem
	err = db.AutoMigrate(
		&types.Material{},
		&types.Machine{},
		&types.ProductionOrder{},
		&types.OrderItem{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}

	// Inject test DB
	storage.OverrideDatabaseInstance(db)
	return db
}

// Helper to quickly seed dependencies
func seedDependencies(db *gorm.DB) (int, int) {
	mat := types.Material{Name: "PLA Test", Type: "FDM", IdCompany: testCompanyIDOrder}
	mac := types.Machine{Name: "Prusa Test", Type: "FDM", IdCompany: testCompanyIDOrder}

	db.Create(&mat)
	db.Create(&mac)

	return mat.ID, mac.ID
}

func TestUpdateOrderUseCase(t *testing.T) {

	t.Run("Update Basic Data and Progress", func(t *testing.T) {
		db := setupOrderTest(t)
		matID, macID := seedDependencies(db)

		// 1. Create an initial order with explicit ID
		initialOrder := types.ProductionOrder{
			IdOrder:          1,
			IdCompany:        testCompanyIDOrder,
			ClientName:       "Old Client",
			TotalPieces:      10,
			Status:           "pending",
			OperatorID:       1,
			EstimatedMinutes: 60,
			Items: []types.OrderItem{
				{StlName: "Part A", Quantity: 5, MaterialID: &matID, MachineID: &macID},
			},
		}
		db.Create(&initialOrder)

		// 2. Prepare update DTO
		orderID := uint(1)
		clientName := "New Client"
		totalPieces := 10
		donePieces := 5
		status := "in-progress"
		priority := "P1"
		notes := "Making fast progress"
		operatorID := 1
		estimatedMinutes := 120
		deadlineString := "2025-12-31"

		updateDTO := types.UpdateOrderDTO{
			ID:               &orderID,
			ClientName:       &clientName,
			TotalPieces:      &totalPieces,
			DonePieces:       &donePieces,
			Status:           &status,
			Priority:         &priority,
			Notes:            &notes,
			OperatorID:       &operatorID,
			EstimatedMinutes: &estimatedMinutes,
			Deadline:         &deadlineString,
		}

		// 3. Execute Update
		err := UpdateOrderUseCase(int(initialOrder.IdOrder), updateDTO, testCompanyIDOrder)
		if err != nil {
			t.Fatalf("Failed to update order: %v", err)
		}

		// 4. Verify Changes
		var order types.ProductionOrder
		db.Preload("Items").First(&order, initialOrder.IdOrder)

		if order.ClientName != "New Client" {
			t.Errorf("Client name not updated. Got: %s", order.ClientName)
		}
		if order.DonePieces != 5 {
			t.Errorf("Progress not updated. Got: %d", order.DonePieces)
		}
		if order.EstimatedMinutes != 120 {
			t.Errorf("Time calculation error. Expected 120, Got: %d", order.EstimatedMinutes)
		}

		// Verify deadline
		expectedDate, _ := time.Parse("2006-01-02", "2025-12-31")
		if !order.Deadline.Equal(expectedDate) {
			t.Error("Deadline was not saved correctly")
		}
	})

	t.Run("Auto-Complete Order (Business Logic)", func(t *testing.T) {
		db := setupOrderTest(t)
		matID, macID := seedDependencies(db)

		// Create order with explicit ID (required by autoIncrement:false)
		order := types.ProductionOrder{
			IdOrder:     100,
			IdCompany:   testCompanyIDOrder,
			ClientName:  "Test Client",
			TotalPieces: 10,
			DonePieces:  0,
			Status:      "in-progress",
			OperatorID:  1,
			Deadline:    time.Now().Add(24 * time.Hour),
		}
		db.Create(&order)

		// Send items with all pieces completed
		items := []types.CreateOrderItemDTO{
			{ID: 1, StlName: "Part A", Quantity: 5, DonePieces: 5, MaterialID: &matID, MachineID: &macID},
			{ID: 2, StlName: "Part B", Quantity: 5, DonePieces: 5, MaterialID: &matID},
		}

		status := "in-progress"
		operatorID := 1

		dto := types.UpdateOrderDTO{
			OperatorID: &operatorID,
			Status:     &status,
			Items:      &items,
		}

		err := UpdateOrderUseCase(int(order.IdOrder), dto, testCompanyIDOrder)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		var updatedOrder types.ProductionOrder
		db.First(&updatedOrder, order.IdOrder)

		// The UseCase detects DonePieces >= TotalPieces and changes Status to "ready"
		if updatedOrder.Status != "ready" {
			t.Errorf("Order should have been auto-completed to 'ready'. Current status: %s", updatedOrder.Status)
		}
	})

	t.Run("Update Non-Existent Order", func(t *testing.T) {
		setupOrderTest(t) // Empty DB

		err := UpdateOrderUseCase(999, types.UpdateOrderDTO{}, testCompanyIDOrder)
		if err == nil {
			t.Error("Expected error because order 999 does not exist")
		}
	})
}
