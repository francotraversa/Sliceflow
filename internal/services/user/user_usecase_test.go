package services

import (
	"testing"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error al abrir DB de test: %v", err)
	}
	db.AutoMigrate(&types.User{})

	storage.OverrideDatabaseInstance(db)
	return db
}

func TestDeleteUserUseCase(t *testing.T) {
	db := setupTest(t)

	user := types.User{Username: "borrame", Role: "user", Status: "active"}
	db.Create(&user)

	t.Run("Usuario se borra a sí mismo (éxito)", func(t *testing.T) {
		err := DeleteUserUseCase(user.ID, user.ID, "user")
		if err != nil {
			t.Errorf("No debería haber error: %v", err)
		}

		var found types.User
		db.First(&found, user.ID)
		if found.Status != "disabled" {
			t.Errorf("Se esperaba status disabled, se obtuvo %s", found.Status)
		}
	})

	t.Run("Usuario intenta borrar a otro (error)", func(t *testing.T) {
		err := DeleteUserUseCase(user.ID, 999, "user")
		if err == nil {
			t.Error("Se esperaba error de permisos")
		}
	})
}

func TestUpdateUserUseCase(t *testing.T) {
	db := setupTest(t)

	u := types.User{Username: "antiguo", Role: "user", Status: "active"}
	db.Create(&u)

	t.Run("Cambio de Username exitoso", func(t *testing.T) {
		update := types.UserUpdateCreds{Username: "nuevo"}
		err := UpdateUserUseCase(u.ID, u.ID, "user", update)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		var found types.User
		db.First(&found, u.ID)
		if found.Username != "nuevo" {
			t.Errorf("Se esperaba 'nuevo', se obtuvo %s", found.Username)
		}
	})

	t.Run("Usuario intenta cambiar su Rol (denegado)", func(t *testing.T) {
		update := types.UserUpdateCreds{Role: "admin"}
		err := UpdateUserUseCase(u.ID, u.ID, "user", update)
		if err == nil {
			t.Error("Se esperaba error de restricción de admin")
		}
	})
}

func TestGetAllUserUseCase(t *testing.T) {
	db := setupTest(t)

	db.Create(&types.User{Username: "admin1", Role: "admin", Status: "active"})
	db.Create(&types.User{Username: "user1", Role: "user", Status: "active"})
	db.Create(&types.User{Username: "user2", Role: "user", Status: "disabled"})

	t.Run("Listar solo Admins", func(t *testing.T) {
		users, err := GetAllUserUserUseCase("admin", "admin")
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if len(*users) != 1 || (*users)[0].Username != "admin1" {
			t.Errorf("Filtro de admin falló. Obtenidos: %d", len(*users))
		}
	})

	t.Run("Listar Todos (sin filtro)", func(t *testing.T) {
		users, err := GetAllUserUserUseCase("admin", "")
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if len(*users) != 3 {
			t.Errorf("Se esperaban 3 usuarios, se obtuvieron %d", len(*users))
		}
	})
}
