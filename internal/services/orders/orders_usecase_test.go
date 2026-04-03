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

	err = db.AutoMigrate(
		&types.Material{},
		&types.Machine{},
		&types.ProductionOrder{},
		&types.OrderItem{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}

	storage.OverrideDatabaseInstance(db)
	return db
}

// seedDependencies: Creates material and machine for use in tests
func seedDependencies(db *gorm.DB) (int, int) {
	mat := types.Material{Name: "PLA Test", Type: "FDM", IdCompany: testCompanyIDOrder}
	mac := types.Machine{Name: "Prusa Test", Type: "FDM", IdCompany: testCompanyIDOrder}

	db.Create(&mat)
	db.Create(&mac)

	return mat.ID, mac.ID
}

// ptr helpers
func ptrFloat(v float64) *float64 { return &v }
func ptrInt(v int) *int           { return &v }

// ─────────────────────────────────────────────
// CREATE ORDER TESTS
// ─────────────────────────────────────────────

func TestCreateOrderUseCase(t *testing.T) {

	t.Run("Create Order with Weight and Time on Items", func(t *testing.T) {
		db := setupOrderTest(t)
		matID, macID := seedDependencies(db)

		id := uint(1)
		dto := types.CreateOrderDTO{
			ID:               &id,
			ClientName:       "Test Client",
			Priority:         "P1",
			Notes:            "test notes",
			EstimatedMinutes: 30,
			Deadline:         "2025-12-31",
			OperatorID:       1,
			Items: []types.CreateOrderItemDTO{
				{
					StlName:    "Part A",
					Quantity:   3,
					MaterialID: &matID,
					MachineID:  &macID,
					Price:      15.0,
					Weight:     ptrFloat(120.5),
					Time:       ptrInt(45),
				},
				{
					StlName:    "Part B",
					Quantity:   2,
					MaterialID: &matID,
					Price:      10.0,
					Weight:     ptrFloat(80.0),
					// Time omitted → nil
				},
			},
		}

		err := CreateOrderUseCase(dto, testCompanyIDOrder)
		if err != nil {
			t.Fatalf("CreateOrderUseCase failed: %v", err)
		}

		var order types.ProductionOrder
		db.Preload("Items").Where("id_order = ?", id).First(&order)

		if len(order.Items) != 2 {
			t.Fatalf("Expected 2 items, got %d", len(order.Items))
		}

		// Verify Weight and Time on first item
		item0 := order.Items[0]
		if item0.Weight == nil || *item0.Weight != 120.5 {
			t.Errorf("Item 0 Weight: expected 120.5, got %v", item0.Weight)
		}
		if item0.Time == nil || *item0.Time != 45 {
			t.Errorf("Item 0 Time: expected 45, got %v", item0.Time)
		}

		// Verify Weight set, Time nil on second item
		item1 := order.Items[1]
		if item1.Weight == nil || *item1.Weight != 80.0 {
			t.Errorf("Item 1 Weight: expected 80.0, got %v", item1.Weight)
		}
		if item1.Time != nil {
			t.Errorf("Item 1 Time: expected nil, got %v", item1.Time)
		}

		// Verify total pieces and price
		if order.TotalPieces != 5 {
			t.Errorf("TotalPieces: expected 5, got %d", order.TotalPieces)
		}
	})

	t.Run("Create Order without Weight and Time (backward compat)", func(t *testing.T) {
		db := setupOrderTest(t)
		matID, _ := seedDependencies(db)

		id := uint(2)
		dto := types.CreateOrderDTO{
			ID:         &id,
			ClientName: "Simple Client",
			Priority:   "P3",
			Deadline:   "2025-06-30",
			OperatorID: 1,
			Items: []types.CreateOrderItemDTO{
				{StlName: "Part X", Quantity: 1, MaterialID: &matID, Price: 5.0},
			},
		}

		err := CreateOrderUseCase(dto, testCompanyIDOrder)
		if err != nil {
			t.Fatalf("CreateOrderUseCase failed: %v", err)
		}

		var order types.ProductionOrder
		db.Preload("Items").Where("id_order = ?", id).First(&order)

		if len(order.Items) != 1 {
			t.Fatalf("Expected 1 item, got %d", len(order.Items))
		}
		if order.Items[0].Weight != nil || order.Items[0].Time != nil {
			t.Error("Weight and Time should be nil when not provided")
		}
	})

	t.Run("Create Duplicate Order returns error", func(t *testing.T) {
		db := setupOrderTest(t)
		matID, _ := seedDependencies(db)

		id := uint(99)
		dto := types.CreateOrderDTO{
			ID:         &id,
			ClientName: "Client",
			Priority:   "P3",
			Deadline:   "2025-12-31",
			OperatorID: 1,
			Items:      []types.CreateOrderItemDTO{{StlName: "X", Quantity: 1, MaterialID: &matID}},
		}

		_ = CreateOrderUseCase(dto, testCompanyIDOrder)
		err := CreateOrderUseCase(dto, testCompanyIDOrder) // duplicate
		if err == nil {
			t.Error("Expected error when creating duplicate order")
		}
	})
}

// ─────────────────────────────────────────────
// UPDATE ORDER TESTS
// ─────────────────────────────────────────────

func TestUpdateOrderUseCase(t *testing.T) {

	t.Run("Update Basic Data and Progress", func(t *testing.T) {
		db := setupOrderTest(t)
		matID, macID := seedDependencies(db)

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

		err := UpdateOrderUseCase(int(initialOrder.Id), updateDTO, testCompanyIDOrder)
		if err != nil {
			t.Fatalf("Failed to update order: %v", err)
		}

		var order types.ProductionOrder
		db.Preload("Items").Where("id = ?", initialOrder.Id).First(&order)

		if order.ClientName != "New Client" {
			t.Errorf("ClientName: expected 'New Client', got '%s'", order.ClientName)
		}
		if order.DonePieces != 5 {
			t.Errorf("DonePieces: expected 5, got %d", order.DonePieces)
		}
		if order.EstimatedMinutes != 120 {
			t.Errorf("EstimatedMinutes: expected 120, got %d", order.EstimatedMinutes)
		}

		expectedDate, _ := time.Parse("2006-01-02", "2025-12-31")
		if !order.Deadline.Equal(expectedDate) {
			t.Error("Deadline was not saved correctly")
		}
	})

	t.Run("Update Items with Weight and Time", func(t *testing.T) {
		db := setupOrderTest(t)
		matID, macID := seedDependencies(db)

		order := types.ProductionOrder{
			IdOrder:     10,
			IdCompany:   testCompanyIDOrder,
			ClientName:  "Client",
			TotalPieces: 5,
			Status:      "pending",
			OperatorID:  1,
			Deadline:    time.Now().Add(24 * time.Hour),
			Items: []types.OrderItem{
				{StlName: "Old Part", Quantity: 5},
			},
		}
		db.Create(&order)

		items := []types.CreateOrderItemDTO{
			{
				StlName:    "New Part A",
				Quantity:   3,
				DonePieces: 1,
				MaterialID: &matID,
				MachineID:  &macID,
				Price:      20.0,
				Weight:     ptrFloat(200.0),
				Time:       ptrInt(90),
			},
		}
		operatorID := 1
		dto := types.UpdateOrderDTO{
			OperatorID: &operatorID,
			Items:      &items,
		}

		err := UpdateOrderUseCase(int(order.Id), dto, testCompanyIDOrder)
		if err != nil {
			t.Fatalf("UpdateOrderUseCase failed: %v", err)
		}

		var updated types.ProductionOrder
		db.Preload("Items").Where("id = ?", order.Id).First(&updated)

		if len(updated.Items) != 1 {
			t.Fatalf("Expected 1 item after update, got %d", len(updated.Items))
		}
		item := updated.Items[0]
		if item.StlName != "New Part A" {
			t.Errorf("StlName: expected 'New Part A', got '%s'", item.StlName)
		}
		if item.Weight == nil || *item.Weight != 200.0 {
			t.Errorf("Weight: expected 200.0, got %v", item.Weight)
		}
		if item.Time == nil || *item.Time != 90 {
			t.Errorf("Time: expected 90, got %v", item.Time)
		}
	})

	t.Run("Auto-Complete Order (Business Logic)", func(t *testing.T) {
		db := setupOrderTest(t)
		matID, macID := seedDependencies(db)

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

		items := []types.CreateOrderItemDTO{
			{StlName: "Part A", Quantity: 5, DonePieces: 5, MaterialID: &matID, MachineID: &macID},
			{StlName: "Part B", Quantity: 5, DonePieces: 5, MaterialID: &matID},
		}
		status := "in-progress"
		operatorID := 1

		dto := types.UpdateOrderDTO{
			OperatorID: &operatorID,
			Status:     &status,
			Items:      &items,
		}

		err := UpdateOrderUseCase(int(order.Id), dto, testCompanyIDOrder)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		var updatedOrder types.ProductionOrder
		db.Where("id = ?", order.Id).First(&updatedOrder)

		// DonePieces(10) >= TotalPieces(10) → status should be "ready"
		if updatedOrder.Status != "ready" {
			t.Errorf("Expected status 'ready', got '%s'", updatedOrder.Status)
		}
	})

	t.Run("Update Non-Existent Order returns error", func(t *testing.T) {
		setupOrderTest(t)

		err := UpdateOrderUseCase(999, types.UpdateOrderDTO{}, testCompanyIDOrder)
		if err == nil {
			t.Error("Expected error for non-existent order 999")
		}
	})
}
