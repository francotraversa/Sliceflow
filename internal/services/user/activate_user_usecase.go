package services

import (
	"fmt"

	userStorage "github.com/francotraversa/Sliceflow/internal/infra/database/user_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func EnableUserByIDUseCase(req types.UserIDActivate) error {
	if req.IdUser == 0 {
		return fmt.Errorf("invalid user ID")
	}
	user := userStorage.FindUserByUserId(req.IdUser)
	if user == nil {
		return fmt.Errorf("user not found with ID %d", req.IdUser)
	}

	if user.Status == "active" {
		return fmt.Errorf("user is already active")
	}

	user.Status = "active"

	if err := db_utils.SaveWithoutCompany(user); err != nil {
		return fmt.Errorf("failed to update user in database")
	}

	return nil
}
