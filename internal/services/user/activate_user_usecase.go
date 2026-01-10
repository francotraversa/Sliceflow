package services

import (
	"fmt"

	userStorage "github.com/francotraversa/Sliceflow/internal/infra/database/user_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func EnableUserByIDUseCase(req types.UserIDActivate) error {
	// 1. Validar que el ID no sea 0
	if req.ID == 0 {
		return fmt.Errorf("invalid user ID")
	}

	// 2. Buscar usuario por ID
	user := userStorage.FindUserByUserId(req.ID)
	if user == nil {
		return fmt.Errorf("user not found with ID %d", req.ID)
	}

	if user.Status == "active" {
		return fmt.Errorf("user is already active")
	}

	user.Status = "active"

	if err := db_utils.Save(user); err != nil {
		return fmt.Errorf("failed to update user in database")
	}

	return nil
}
