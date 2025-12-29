package services

import (
	"testing"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupMachineTest crea una DB limpia para cada prueba
func setupMachineTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error al iniciar DB de test: %v", err)
	}

	// Migramos la tabla Machine
	err = db.AutoMigrate(&types.Machine{})
	if err != nil {
		t.Fatalf("Error al migrar: %v", err)
	}

	// Inyectamos la DB en la instancia global
	// Asignamos directamente a la variable global DB
	storage.DBInstance.DB = db
	return db
}

func TestMachineCRUD(t *testing.T) {

	// --- 1. TEST DE CREACIÓN ---
	t.Run("Crear Maquina Exitosa", func(t *testing.T) {
		db := setupMachineTest(t)

		dto := types.CreateMachineDTO{
			Name: "Prusa MK3S - Sector A",
			Type: "FDM",
		}

		err := CreateMachineUseCase(dto)
		if err != nil {
			t.Fatalf("Error inesperado al crear: %v", err)
		}

		// Verificar que se guardó y el estado default es 'idle'
		var machine types.Machine
		db.First(&machine, "name = ?", dto.Name)

		if machine.ID == 0 {
			t.Error("No se guardó la máquina en la DB")
		}
		if machine.Status != "idle" {
			t.Errorf("El estado inicial debió ser 'idle', fue: %s", machine.Status)
		}
	})

	t.Run("Crear Maquina Duplicada (Debe fallar)", func(t *testing.T) {
		setupMachineTest(t)

		dto := types.CreateMachineDTO{Name: "Ender 3 V2", Type: "FDM"}

		// Primera vez
		CreateMachineUseCase(dto)

		// Segunda vez (Mismo nombre)
		err := CreateMachineUseCase(dto)
		if err == nil {
			t.Error("Debió fallar por nombre duplicado, pero pasó")
		}
	})

	// --- 2. TEST DE LECTURA ---
	t.Run("Listar Maquinas", func(t *testing.T) {
		setupMachineTest(t)

		CreateMachineUseCase(types.CreateMachineDTO{Name: "M1", Type: "FDM"})
		CreateMachineUseCase(types.CreateMachineDTO{Name: "M2", Type: "SLS"})

		// CORRECCIÓN: Ahora GetAll recibe un filtro. Pasamos vacío para traer todas.
		list, err := GetAllMachinesUseCase(types.MachineFilter{})
		if err != nil {
			t.Fatalf("Error al listar: %v", err)
		}

		// CORRECCIÓN: 'list' es un slice, no un puntero. Quitamos el *.
		if len(*list) != 2 {
			t.Errorf("Esperaba 2 máquinas, obtuvo %d", len(*list))
		}
	})

	// --- 3. TEST DE ACTUALIZACIÓN ---
	t.Run("Actualizar Maquina", func(t *testing.T) {
		db := setupMachineTest(t)

		// Crear original
		CreateMachineUseCase(types.CreateMachineDTO{Name: "Original"})

		// Update
		id := 1
		updateDTO := types.UpdateMachineDTO{
			Name:   "Editada",
			Type:   "SLA",
			Status: "maintenance", // Probamos cambiar estado manualmente
		}

		err := UpdateMachineUseCase(id, updateDTO)
		if err != nil {
			t.Fatalf("Error al actualizar: %v", err)
		}

		// Verificar
		var m types.Machine
		db.First(&m, id)

		if m.Name != "Editada" || m.Status != "maintenance" {
			t.Errorf("Datos incorrectos. Name: %s, Status: %s", m.Name, m.Status)
		}
	})

	// --- 4. TEST DE BORRADO (SOFT DELETE) ---
	t.Run("Eliminar Maquina (Soft Delete)", func(t *testing.T) {
		db := setupMachineTest(t)

		// Crear
		CreateMachineUseCase(types.CreateMachineDTO{Name: "Para Borrar", Type: "FDM"})

		// Borrar ID 1
		err := DeleteMachineUseCase(1)
		if err != nil {
			t.Fatalf("Error al borrar: %v", err)
		}

		// A. Verificar que GetAll NO la trae
		// CORRECCIÓN: Pasamos filtro vacío
		list, _ := GetAllMachinesUseCase(types.MachineFilter{})

		// CORRECCIÓN: Quitamos el *
		if len(*list) != 0 {
			t.Errorf("El listado debió venir vacío, trajo %d", len(*list))
		}

		// B. Verificar que sigue en DB con deleted_at (Unscoped)
		var count int64
		db.Unscoped().Model(&types.Machine{}).Where("id = ?", 1).Count(&count)
		if count != 1 {
			t.Error("El registro físico desapareció (Debió ser borrado lógico)")
		}
	})
}
