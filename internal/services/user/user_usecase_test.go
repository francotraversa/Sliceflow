package services

import (
	"testing"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testCompanyIDUser int = 1

func setupTest(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}
	db.AutoMigrate(&types.User{})

	storage.OverrideDatabaseInstance(db)
	return db
}

func TestDeleteUserUseCase(t *testing.T) {
	db := setupTest(t)

	user := types.User{Username: "deleteme", Role: "user", Status: "active"}
	db.Create(&user)

	t.Run("User deletes themselves (success)", func(t *testing.T) {
		err := DeleteUserUseCase(user.IdUser, user.IdUser, "user")
		if err != nil {
			t.Errorf("Should not have errored: %v", err)
		}

		var found types.User
		db.First(&found, user.IdUser)
		if found.Status != "disabled" {
			t.Errorf("Expected status disabled, got %s", found.Status)
		}
	})

	t.Run("User tries to delete another user (error)", func(t *testing.T) {
		err := DeleteUserUseCase(user.IdUser, 999, "user")
		if err == nil {
			t.Error("Expected permission error")
		}
	})
}

func TestUpdateUserUseCase(t *testing.T) {
	db := setupTest(t)

	u := types.User{Username: "oldname", Role: "user", Status: "active"}
	db.Create(&u)

	t.Run("Username Change Successful", func(t *testing.T) {
		update := types.UserUpdateCreds{Username: "newname"}
		err := UpdateUserUseCase(u.IdUser, u.IdUser, "user", update)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		var found types.User
		db.First(&found, u.IdUser)
		if found.Username != "newname" {
			t.Errorf("Expected 'newname', got %s", found.Username)
		}
	})

	t.Run("User tries to change own role (denied)", func(t *testing.T) {
		update := types.UserUpdateCreds{Role: "admin"}
		err := UpdateUserUseCase(u.IdUser, u.IdUser, "user", update)
		if err == nil {
			t.Error("Expected admin restriction error")
		}
	})
}

func TestGetAllUserUseCase(t *testing.T) {
	db := setupTest(t)

	db.Create(&types.User{Username: "admin1", Role: "admin", Status: "active", IdCompany: uint(testCompanyIDUser)})
	db.Create(&types.User{Username: "user1", Role: "user", Status: "active", IdCompany: uint(testCompanyIDUser)})
	db.Create(&types.User{Username: "user2", Role: "user", Status: "disabled", IdCompany: uint(testCompanyIDUser)})

	t.Run("List Admins Only", func(t *testing.T) {
		users, err := GetAllUserUserUseCase("admin", "admin", "admin", testCompanyIDUser)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if len(*users) != 1 || (*users)[0].Username != "admin1" {
			t.Errorf("Admin filter failed. Got: %d", len(*users))
		}
	})

	t.Run("List All (no filter)", func(t *testing.T) {
		users, err := GetAllUserUserUseCase("admin", "", "admin", testCompanyIDUser)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if len(*users) != 3 {
			t.Errorf("Expected 3 users, got %d", len(*users))
		}
	})
}
