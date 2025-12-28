package services

import (
	"testing"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupStockTest adapta tu función de setup para el catálogo de productos
func setupStockTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error al abrir DB de test: %v", err)
	}

	// Migramos el modelo de Stock
	db.AutoMigrate(&types.StockItem{})

	// Inyectamos la base de datos de test en tu instancia global
	storage.OverrideDatabaseInstance(db)
	return db
}
func TestStockService(t *testing.T) {
	db := setupStockTest(t)

	t.Run("Crear Producto Exitoso", func(t *testing.T) {
		req := types.ProductCreateRequest{
			SKU:         "SKU-779-1",
			Name:        "Tornillo Phillips",
			Description: "1/2 pulgada",
			Quantity:    100,
		}

		err := CreateProductUseCase(req)
		if err != nil {
			t.Fatalf("No debería haber error al crear producto: %v", err)
		}

		var found types.StockItem
		db.Where("sku = ?", "SKU-779-1").First(&found)
		if found.Name != "Tornillo Phillips" {
			t.Errorf("Se esperaba 'Tornillo Phillips', se obtuvo %s", found.Name)
		}
	})

	t.Run("Evitar SKU Duplicado", func(t *testing.T) {
		req := types.ProductCreateRequest{
			SKU:  "SKU-779-1", // SKU que ya existe del test anterior
			Name: "Producto Duplicado",
		}

		err := CreateProductUseCase(req)
		if err == nil {
			t.Error("Se esperaba un error por SKU duplicado")
		}
	})

	t.Run("Soft Delete de Producto", func(t *testing.T) {
		// Creamos uno nuevo para borrar
		item := types.StockItem{SKU: "DEL-99", Name: "Borrame"}
		db.Create(&item)

		err := DeleteByIdUseCase(item.ID)
		if err != nil {
			t.Fatalf("Error al borrar: %v", err)
		}

		// 1. Verificar que First() no lo encuentra (Consulta Normal)
		var found types.StockItem
		result := db.First(&found, item.ID)
		if result.Error == nil {
			t.Error("El producto sigue visible en consultas normales")
		}

		// 2. Verificar que sigue en la DB físicamente (Unscoped)
		var raw types.StockItem
		db.Unscoped().First(&raw, item.ID)
		if !raw.DeletedAt.Valid {
			t.Error("DeletedAt no fue completado por GORM")
		}
	})

	t.Run("Actualizar Producto - Exito y Fallo", func(t *testing.T) {
		// A. Caso Feliz
		item := types.StockItem{SKU: "EDIT-1", Name: "Viejo", Quantity: 50}
		db.Create(&item)

		updateReq := types.ProductUpdateRequest{Name: "Nuevo Nombre"}
		_, err := UpdateByIdProductUseCase(item.SKU, updateReq)
		if err != nil {
			t.Errorf("Error inesperado al actualizar: %v", err)
		}

		// Validamos cambio
		var found types.StockItem
		db.First(&found, item.ID)
		if found.Name != "Nuevo Nombre" {
			t.Errorf("No se actualizó el nombre. Valor actual: %s", found.Name)
		}

		// B. Caso Fallido (ID no existe)
		_, err = UpdateByIdProductUseCase("99999", updateReq)
		if err == nil {
			t.Error("Se esperaba error al actualizar ID inexistente, pero no falló")
		}
	})
}
