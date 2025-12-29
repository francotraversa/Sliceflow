package services

import (
	"testing"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupMaterialTest crea una DB en memoria limpia para cada prueba
func setupMaterialTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error al iniciar DB de test: %v", err)
	}

	// Migramos solo la tabla de Materiales
	err = db.AutoMigrate(&types.Material{})
	if err != nil {
		t.Fatalf("Error al migrar: %v", err)
	}

	// Inyectamos la DB de prueba en la instancia global
	// CORRECCIÓN: Asignación directa
	storage.DBInstance.DB = db
	return db
}

func TestMaterialCRUD(t *testing.T) {
	// --- 1. TEST DE CREACIÓN ---
	t.Run("Crear Material Exitoso", func(t *testing.T) {
		db := setupMaterialTest(t)

		dto := types.CreateMaterialDTO{
			Name:        "PLA Negro Grilon",
			Type:        "Filamento",
			Description: "Bobina 1kg",
		}

		err := CreateMaterialUseCase(dto)
		if err != nil {
			t.Fatalf("Error inesperado al crear: %v", err)
		}

		// Verificar que se guardó
		var count int64
		db.Model(&types.Material{}).Count(&count)
		if count != 1 {
			t.Errorf("Esperaba 1 material, hay %d", count)
		}
	})

	t.Run("Crear Material Duplicado (Debe fallar)", func(t *testing.T) {
		setupMaterialTest(t)

		dto := types.CreateMaterialDTO{Name: "PLA Unico", Type: "Filamento"}

		// Primer insert
		_ = CreateMaterialUseCase(dto)

		// Segundo insert (Mismo nombre)
		err := CreateMaterialUseCase(dto)
		if err == nil {
			t.Error("Debió fallar por nombre duplicado, pero pasó")
		}
	})

	// --- 2. TEST DE LECTURA (LISTADO) ---
	t.Run("Listar Materiales", func(t *testing.T) {
		setupMaterialTest(t)

		// Insertamos 2 materiales
		CreateMaterialUseCase(types.CreateMaterialDTO{Name: "Mat 1", Type: "A"})
		CreateMaterialUseCase(types.CreateMaterialDTO{Name: "Mat 2", Type: "B"})

		// CORRECCIÓN: Pasamos el filtro vacío types.MaterialFilter{}
		result, err := GetAllMaterialsUseCase(types.MaterialFilter{})
		if err != nil {
			t.Fatalf("Error al listar: %v", err)
		}

		// CORRECCIÓN: Quitamos el *, es un slice normal
		if len(*result) != 2 {
			t.Errorf("Esperaba 2 materiales, obtuvo %d", len(*result))
		}
	})

	// --- 3. TEST DE ACTUALIZACIÓN ---
	t.Run("Actualizar Material", func(t *testing.T) {
		db := setupMaterialTest(t)

		CreateMaterialUseCase(types.CreateMaterialDTO{Name: "Nombre Viejo", Type: "Viejo"})

		id := 1
		updateDTO := types.UpdateMaterialDTO{
			Name:        "Nombre Nuevo",
			Type:        "Nuevo Tipo",
			Description: "Editado",
		}

		err := UpdateMaterialUseCase(id, updateDTO)
		if err != nil {
			t.Fatalf("Error al actualizar: %v", err)
		}

		var mat types.Material
		db.First(&mat, id)
		if mat.Name != "Nombre Nuevo" {
			t.Errorf("No se actualizó el nombre. Valor: %s", mat.Name)
		}
	})

	// --- 4. TEST DE BORRADO (SOFT DELETE) ---
	t.Run("Eliminar Material (Soft Delete)", func(t *testing.T) {
		db := setupMaterialTest(t)

		// Crear material ID 1
		CreateMaterialUseCase(types.CreateMaterialDTO{Name: "A Borrar", Type: "X"})

		// Eliminar
		err := DeleteMaterialUseCase(1)
		if err != nil {
			t.Fatalf("Error al borrar: %v", err)
		}

		// Verificaciones:

		// A. GetAll no lo debe traer
		// CORRECCIÓN: Filtro vacío
		list, _ := GetAllMaterialsUseCase(types.MaterialFilter{})

		// CORRECCIÓN: Quitamos el *
		if len(*list) != 0 {
			t.Errorf("El listado debió venir vacío, trajo %d", len(*list))
		}

		// B. En la DB física el registro SIGUE existiendo (Soft Delete)
		var count int64
		db.Unscoped().Model(&types.Material{}).Where("id = ?", 1).Count(&count)
		if count != 1 {
			t.Error("El registro físico desapareció de la DB (Debió ser borrado lógico)")
		}

		// C. Verificar que deleted_at no sea nulo
		var mat types.Material
		db.Unscoped().First(&mat, 1)
		if !mat.DeletedAt.Valid {
			t.Error("El campo DeletedAt está vacío")
		}
	})
}
