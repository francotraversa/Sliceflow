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
			Price:       90,
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
		item := types.StockItem{SKU: "DEL-99", Name: "Borrame", Status: "active"}
		db.Create(&item)

		// Llamamos al UseCase (Asumiendo que ya lo actualizaste para recibir string)
		err := DeleteByIdUseCase(item.SKU)
		if err != nil {
			t.Fatalf("Error al borrar: %v", err)
		}

		// 1. Verificar que NO lo encuentra en consulta normal
		var found types.StockItem
		// CORRECCIÓN: Usamos Where explícito porque SKU es un string
		result := db.Where("sku = ?", item.SKU).First(&found)

		if result.Error == nil {
			t.Error("El producto sigue visible en consultas normales")
		}

		// 2. Verificar que SIGUE en la DB físicamente (Unscoped)
		var raw types.StockItem
		// CORRECCIÓN: Usamos Where explícito también aquí
		db.Unscoped().Where("sku = ?", item.SKU).First(&raw)

		if !raw.DeletedAt.Valid {
			t.Error("DeletedAt no fue completado por GORM")
		}
	})

	t.Run("Actualizar Producto - Exito y Fallo", func(t *testing.T) {
		// A. Caso Feliz
		item := types.StockItem{SKU: "EDIT-1", Name: "Viejo", Quantity: 50}
		db.Create(&item)

		updateReq := types.ProductUpdateRequest{Name: "Nuevo Nombre"}

		// Llamada al UseCase con SKU string
		_, err := UpdateByIdProductUseCase(item.SKU, updateReq)
		if err != nil {
			t.Errorf("Error inesperado al actualizar: %v", err)
		}

		// Validamos cambio
		var found types.StockItem
		// CORRECCIÓN: Búsqueda explícita por SKU
		db.Where("sku = ?", item.SKU).First(&found)

		if found.Name != "Nuevo Nombre" {
			t.Errorf("No se actualizó el nombre. Valor actual: %s", found.Name)
		}

		// B. Caso Fallido (SKU no existe)
		_, err = UpdateByIdProductUseCase("SKU-FALSO-999", updateReq)
		if err == nil {
			t.Error("Se esperaba error al actualizar SKU inexistente, pero no falló")
		}
	})
}
