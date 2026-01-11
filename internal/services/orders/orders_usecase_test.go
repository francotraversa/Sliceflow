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

	err = db.AutoMigrate(&types.Material{}, &types.Machine{}, &types.ProductionOrder{})
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

		// 1. Crear una Orden Inicial (Manualmente en DB para ir rápido)
		initialOrder := types.ProductionOrder{
			ClientName:       "Cliente Viejo",
			TotalPieces:      10,
			Status:           "pending",
			MaterialID:       matID,
			OperatorID:       1,
			EstimatedMinutes: 60,
		}
		db.Create(&initialOrder)

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

		// Validaciones
		if order.ClientName != "Cliente Nuevo" {
			t.Errorf("No actualizó el cliente. Got: %s", order.ClientName)
		}
		if order.DonePieces != 5 {
			t.Errorf("No actualizó el progreso. Got: %d", order.DonePieces)
		}
		if order.Status != "in-progress" {
			t.Errorf("No actualizó el estado. Got: %s", order.Status)
		}
		if order.EstimatedMinutes != 120 { // 2 horas * 60
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

		// Orden con 10 piezas
		order := types.ProductionOrder{
			TotalPieces: 10,
			DonePieces:  0,
			Status:      "in-progress",
			MaterialID:  matID,
			OperatorID:  1,
		}
		db.Create(&order)

		// Actualizamos diciendo que hizo las 10
		totalPieces := 10
		donePieces := 10
		operatorID := 1
		status := "in-progress"
		dto := types.UpdateOrderDTO{
			TotalPieces: &totalPieces,
			DonePieces:  &donePieces, // <--- TERMINÓ TODO
			MaterialID:  &matID,
			OperatorID:  &operatorID,
			Status:      &status, // El front manda esto, pero el back debería corregirlo
		}

		err := UpdateOrderUseCase(int(order.ID), dto)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		var updatedOrder types.ProductionOrder
		db.First(&updatedOrder, order.ID)

		// El UseCase debería haber detectado Done >= Total y cambiar Status a "completed"
		if updatedOrder.Status != "completed" {
			t.Errorf("La orden debió pasar a 'completed' automáticamente. Estado actual: %s", updatedOrder.Status)
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
