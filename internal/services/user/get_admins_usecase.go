package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllAdminUserUseCase() (*[]types.User, error) {
	var user []types.User
	db := storage.DatabaseInstance{}.Instance()
	if err := db.Where("role = ?", "admin").Find(&user).Error; err != nil {
		return nil, fmt.Errorf("error listing admins")
	}
	return &user, nil
}
