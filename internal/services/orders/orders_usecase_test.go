package services

import (
	"testing"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupOrderTest: Prepara DB en memoria y tablas necesarias
func setupOrderTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error al iniciar DB de test: %v", err)
	}

	// Migrar TODAS las tablas que el use case toca, incluyendo OrderItem
	err = db.AutoMigrate(
		&types.Material{},
		&types.Machine{},
		&types.ProductionOrder{},
		&types.OrderItem{},
	)
	if err != nil {
		t.Fatalf("Error al migrar tablas: %v", err)
	}

	// Inyectar DB
	storage.OverrideDatabaseInstance(db)
	return db
}

// Helper para crear dependencias rápidamente
func seedDependencies(db *gorm.DB) (int, int) {
	mat := types.Material{Name: "PLA Test", Type: "FDM"}
	mac := types.Machine{Name: "Prusa Test", Type: "FDM"}

	db.Create(&mat)
	db.Create(&mac)

	return mat.ID, mac.ID
}

func TestUpdateOrderUseCase(t *testing.T) {

	t.Run("Actualizar Datos Básicos y Progreso", func(t *testing.T) {
		db := setupOrderTest(t)
		matID, macID := seedDependencies(db)

		// 1. Crear una Orden Inicial con ID explícito
		initialOrder := types.ProductionOrder{
			ID:               1,
			ClientName:       "Cliente Viejo",
			TotalPieces:      10,
			Status:           "pending",
			MaterialID:       matID,
			OperatorID:       1,
			EstimatedMinutes: 60,
		}
		db.Create(&initialOrder)

		// 2. Preparar DTO de actualización
		orderID := uint(1)
		clientName := "Cliente Nuevo"
		totalPieces := 10
		donePieces := 5
		status := "in-progress"
		priority := "P1"
		notes := "Avanzando rápido"
		operatorID := 1
		estimatedMinutes := 120
		price := 5000.0
		deadlineString := "2025-12-31"

		updateDTO := types.UpdateOrderDTO{
			ID:               &orderID,
			ClientName:       &clientName,
			TotalPieces:      &totalPieces,
			DonePieces:       &donePieces,
			Status:           &status,
			MaterialID:       &matID,
			Priority:         &priority,
			Notes:            &notes,
			OperatorID:       &operatorID,
			EstimatedMinutes: &estimatedMinutes,
			Price:            &price,
			MachineID:        &macID,
			Deadline:         &deadlineString,
		}

		// 3. Ejecutar Update
		err := UpdateOrderUseCase(int(initialOrder.ID), updateDTO)
		if err != nil {
			t.Fatalf("Error al actualizar orden: %v", err)
		}

		// 4. Verificar Cambios
		var order types.ProductionOrder
		db.Preload("Machine").First(&order, initialOrder.ID)

		if order.ClientName != "Cliente Nuevo" {
			t.Errorf("No actualizó el cliente. Got: %s", order.ClientName)
		}
		if order.DonePieces != 5 {
			t.Errorf("No actualizó el progreso. Got: %d", order.DonePieces)
		}
		if order.EstimatedMinutes != 120 {
			t.Errorf("Error calculando tiempo. Esperaba 120, Got: %d", order.EstimatedMinutes)
		}
		if order.MachineID == nil || *order.MachineID != macID {
			t.Error("No se asignó la máquina correctamente")
		}

		// Verificar fecha
		expectedDate, _ := time.Parse("2006-01-02", "2025-12-31")
		if !order.Deadline.Equal(expectedDate) {
			t.Error("La fecha límite no se guardó bien")
		}
	})

	t.Run("Auto-Completar Orden (Lógica de Negocio)", func(t *testing.T) {
		db := setupOrderTest(t)
		matID, _ := seedDependencies(db)

		// Crear orden con ID explícito (necesario por autoIncrement:false)
		order := types.ProductionOrder{
			ID:          100,
			ClientName:  "Test Client",
			TotalPieces: 10,
			DonePieces:  0,
			Status:      "in-progress",
			MaterialID:  matID,
			OperatorID:  1,
			Deadline:    time.Now().Add(24 * time.Hour),
		}
		db.Create(&order)

		// Enviamos items con todas las piezas terminadas
		totalPieces := 10
		donePieces := 10
		operatorID := 1
		status := "in-progress"
		items := []types.CreateOrderItemDTO{
			{ID: 1, ProductName: "Pieza A", Quantity: 5, DonePieces: 5},
			{ID: 2, ProductName: "Pieza B", Quantity: 5, DonePieces: 5},
		}

		dto := types.UpdateOrderDTO{
			TotalPieces: &totalPieces,
			DonePieces:  &donePieces,
			MaterialID:  &matID,
			OperatorID:  &operatorID,
			Status:      &status,
			Items:       &items,
		}

		err := UpdateOrderUseCase(int(order.ID), dto)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		var updatedOrder types.ProductionOrder
		db.First(&updatedOrder, order.ID)

		// El UseCase detecta DonePieces >= TotalPieces y cambia Status a "ready"
		if updatedOrder.Status != "ready" {
			t.Errorf("La orden debió pasar a 'ready' automáticamente. Estado actual: %s", updatedOrder.Status)
		}
	})

	t.Run("Actualizar Orden Inexistente", func(t *testing.T) {
		setupOrderTest(t) // DB vacía

		err := UpdateOrderUseCase(999, types.UpdateOrderDTO{})
		if err == nil {
			t.Error("Debió fallar porque la orden 999 no existe")
		}
	})
}
