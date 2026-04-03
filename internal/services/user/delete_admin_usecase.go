package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	userStorage "github.com/francotraversa/Sliceflow/internal/infra/database/user_utils"
)

func DeleteAdminUseCase(id uint) error {
	db := storage.DatabaseInstance{}.Instance()

	user := userStorage.FindUserByUserId(id)
	if user == nil {
		return fmt.Errorf("The user does not exist")
	}
	if user.Role != "owner" {
		return fmt.Errorf("You are not allowed to delete this user")
	}
	if err := db.Model(&user).Update("status", "disabled").Error; err != nil {
		return fmt.Errorf("Error disabling user")
	}
	return nil
}
