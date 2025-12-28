package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	userStorage "github.com/francotraversa/Sliceflow/internal/database/user_utils"
)

func DeleteUserUseCase(targetID uint, requesterID uint, requesterRole string) error {
	db := storage.DatabaseInstance{}.Instance()

	if requesterID != targetID && requesterRole != "admin" {
		return fmt.Errorf("no tienes permiso para deshabilitar esta cuenta")
	}

	currentUser := userStorage.FindUserByUserId(storage.DBInstance.DB, targetID)
	if currentUser == nil {
		return fmt.Errorf("usuario no encontrado")
	}

	if currentUser.Status == "disabled" {
		return fmt.Errorf("The account is already disable")
	}

	currentUser.Status = "disabled"

	if err := db.Model(&currentUser).Update("status", "disabled").Error; err != nil {
		return fmt.Errorf("Error deactivating a client")
	}

	return nil
}
