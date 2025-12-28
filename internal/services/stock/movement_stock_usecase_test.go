package services

import (
	"testing"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- SETUP & HELPERS ---

// setupMovementTest inicia una DB en memoria limpia para cada test
func setupMovementTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error al abrir DB de test: %v", err)
	}

	// MIGRACIÓN: Agregamos StockItem y StockMovement
	db.AutoMigrate(&types.StockItem{}, &types.StockMovement{})

	storage.OverrideDatabaseInstance(db)
	return db
}

// parseDate ayuda a crear fechas específicas para los tests de filtros
func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}

// --- TEST 1: ESCRITURA (MOVIMIENTOS) ---

func TestStockMovements(t *testing.T) {
	db := setupMovementTest(t)

	// Preparamos un producto base para las pruebas
	baseSku := "TEST-MOV-1"
	initialItem := types.StockItem{
		SKU:      baseSku,
		Name:     "Producto Test",
		Quantity: 10, // Arrancamos con 10 unidades
		Price:    100,
	}
	db.Create(&initialItem)

	t.Run("Movimiento de Entrada (IN) - Suma Stock", func(t *testing.T) {
		req := types.CreateMovementRequest{
			SKU:      baseSku,
			Type:     "IN",
			Quantity: 5, // 10 + 5 = 15
			UserID:   1,
			Reason:   "Compra Proveedor",
		}

		err := AddStockMovementUseCase(req)
		if err != nil {
			t.Fatalf("Error inesperado en movimiento IN: %v", err)
		}

		// 1. Verificar que el StockItem se actualizó
		var item types.StockItem
		db.First(&item, "sku = ?", baseSku)
		if item.Quantity != 15 {
			t.Errorf("El stock debió subir a 15, pero quedó en %d", item.Quantity)
		}

		// 2. Verificar que se creó el historial (Auditoría)
		var mov types.StockMovement
		result := db.Where("stock_sku = ?", baseSku).Last(&mov)
		if result.Error != nil {
			t.Fatalf("No se creó el registro en StockMovements")
		}

		// Validamos la matemática interna
		if mov.QtyBefore != 10 {
			t.Errorf("QtyBefore incorrecto. Esperaba 10, obtuvo %d", mov.QtyBefore)
		}
		if mov.QtyDelta != 5 {
			t.Errorf("QtyDelta incorrecto. Esperaba 5, obtuvo %d", mov.QtyDelta)
		}
		if mov.QtyAfter != 15 {
			t.Errorf("QtyAfter incorrecto. Esperaba 15, obtuvo %d", mov.QtyAfter)
		}
	})

	t.Run("Movimiento de Salida (OUT) - Resta Stock", func(t *testing.T) {
		// Partimos de 15 (resultado del test anterior)
		req := types.CreateMovementRequest{
			SKU:      baseSku,
			Type:     "OUT",
			Quantity: 3, // 15 - 3 = 12
			UserID:   1,
			Reason:   "Venta Cliente",
		}

		err := AddStockMovementUseCase(req)
		if err != nil {
			t.Fatalf("Error inesperado en movimiento OUT: %v", err)
		}

		// 1. Verificar Stock
		var item types.StockItem
		db.First(&item, "sku = ?", baseSku)
		if item.Quantity != 12 {
			t.Errorf("El stock debió bajar a 12, pero quedó en %d", item.Quantity)
		}

		// 2. Verificar Auditoría
		var mov types.StockMovement
		db.Where("stock_sku = ?", baseSku).Last(&mov)

		if mov.QtyBefore != 15 {
			t.Errorf("QtyBefore mal: %d", mov.QtyBefore)
		}
		if mov.QtyDelta != -3 {
			t.Errorf("QtyDelta mal (debe ser negativo): %d", mov.QtyDelta)
		}
		if mov.QtyAfter != 12 {
			t.Errorf("QtyAfter mal: %d", mov.QtyAfter)
		}
	})

	t.Run("Validar Stock Insuficiente (Error)", func(t *testing.T) {
		req := types.CreateMovementRequest{
			SKU:      baseSku,
			Type:     "OUT",
			Quantity: 100, // Imposible (Hay 12)
			UserID:   1,
		}

		err := AddStockMovementUseCase(req)
		if err == nil {
			t.Error("Se esperaba error por falta de stock, pero la operación pasó")
		}

		// Verificar que el stock NO se tocó
		var item types.StockItem
		db.First(&item, "sku = ?", baseSku)
		if item.Quantity != 12 {
			t.Errorf("El stock cambió a %d a pesar del error", item.Quantity)
		}
	})

	t.Run("Fallo por Producto Inexistente", func(t *testing.T) {
		req := types.CreateMovementRequest{
			SKU:      "SKU-FANTASMA",
			Type:     "IN",
			Quantity: 10,
		}

		err := AddStockMovementUseCase(req)
		if err == nil {
			t.Error("Se esperaba error 'producto no encontrado'")
		}
	})
}

// --- TEST 2: LECTURA (HISTORIAL Y FILTROS) ---

func TestGetStockHistory(t *testing.T) {
	db := setupMovementTest(t) // Nueva DB limpia para este test

	// --- SEEDING (Datos de prueba) ---
	skuA := "PROD-A"
	skuB := "PROD-B"

	// Insertamos movimientos manualmente con fechas específicas
	movements := []types.StockMovement{
		// Enero - Producto A
		{StockSKU: skuA, Type: "IN", QtyDelta: 10, CreatedAt: parseDate("2025-01-01")},
		{StockSKU: skuA, Type: "OUT", QtyDelta: -2, CreatedAt: parseDate("2025-01-15")},

		// Febrero - Producto A
		{StockSKU: skuA, Type: "IN", QtyDelta: 5, CreatedAt: parseDate("2025-02-10")},

		// Enero - Producto B
		{StockSKU: skuB, Type: "IN", QtyDelta: 100, CreatedAt: parseDate("2025-01-05")},
	}
	db.Create(&movements)

	// --- EJECUCIÓN ---

	t.Run("Filtrar solo por SKU", func(t *testing.T) {
		// Pedimos solo Producto A
		filter := types.HistoryFilter{SKU: skuA}

		results, err := GetStockHistoryUseCase(filter)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		if len(results) != 3 {
			t.Errorf("Esperaba 3 movimientos para SKU A, obtuvo %d", len(results))
		}
		// Verificar que no venga nada de B
		for _, m := range results {
			if m.StockSKU != skuA {
				t.Errorf("Se filtró mal el SKU: %s", m.StockSKU)
			}
		}
	})

	t.Run("Filtrar solo por Fechas (Enero)", func(t *testing.T) {
		// Pedimos todo lo de Enero (de A y de B)
		filter := types.HistoryFilter{StartDate: "2025-01-01", EndDate: "2025-01-31"}

		results, err := GetStockHistoryUseCase(filter)
		if err != nil {
			t.Fatal(err)
		}

		// Esperamos: 2 de A + 1 de B = 3 Total. (El de Febrero queda fuera)
		if len(results) != 3 {
			t.Errorf("Esperaba 3 movimientos en Enero, obtuvo %d", len(results))
		}
	})

	t.Run("Filtro Combinado (SKU A + Enero)", func(t *testing.T) {
		filter := types.HistoryFilter{
			SKU:       skuA,
			StartDate: "2025-01-01",
			EndDate:   "2025-01-31",
		}

		results, err := GetStockHistoryUseCase(filter)
		if err != nil {
			t.Fatal(err)
		}

		// Esperamos: Solo los 2 de A en Enero.
		if len(results) != 2 {
			t.Errorf("Esperaba 2 movimientos combinados, obtuvo %d", len(results))
		}
	})
}
