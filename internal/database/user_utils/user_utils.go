package userStorage

import (
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func FindUserByUsername(db *gorm.DB, username string) *types.User {
	user := types.User{}
	result := db.Limit(1).Find(&user, "username = ?", username)
	if result.Error == nil && result.RowsAffected > 0 {
		return &user
	}
	return nil
}

func FindUserByUserId(db *gorm.DB, id uint) *types.User {
	var user types.User
	if err := db.First(&user, id).Error; err != nil {
		return nil
	}
	return &user
}

func FindUsersByRole(db *gorm.DB, role string) []types.User {
	var users []types.User
	query := db.Model(&types.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}

	query.Find(&users)
	return users
}
