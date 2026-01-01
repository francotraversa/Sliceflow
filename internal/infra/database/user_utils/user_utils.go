package userStorage

import (
	database "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func FindUserByUsername(username string) *types.User {
	db := database.DatabaseInstance{}.Instance()
	user := types.User{}
	result := db.Limit(1).Find(&user, "username = ?", username)
	if result.Error == nil && result.RowsAffected > 0 {
		return &user
	}
	return nil
}

func FindUserByUserId(id uint) *types.User {
	db := database.DatabaseInstance{}.Instance()
	var user types.User
	if err := db.First(&user, id).Error; err != nil {
		return nil
	}
	return &user
}

func FindUsersByRole(role string) []types.User {
	db := database.DatabaseInstance{}.Instance()
	var users []types.User
	query := db.Model(&types.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}

	query.Find(&users)
	return users
}
