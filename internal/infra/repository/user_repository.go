package repository

import (
	"strings"

	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(id uint) (*types.User, error) {
	var user types.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*types.User, error) {
	var user types.User
	if err := r.db.First(&user, "username = ?", strings.ToLower(username)).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *types.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(id uint, user *types.UserUpdateCreds, companyID uint) error {
	return r.db.Where("id = ? AND id_company = ?", id, companyID).Updates(user).Error
}

func (r *userRepository) DeleteUser(id uint, companyID uint) error {
	return r.db.Delete(&types.User{}, "id = ? AND id_company = ?", id, companyID).Error
}

func (r *userRepository) GetUsers(companyID uint) ([]types.User, error) {
	var users []types.User

	if companyID > 0 {
		if err := r.db.Where("id_company = ?", companyID).Find(&users).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Find(&users).Error; err != nil {
			return nil, err
		}
	}
	return users, nil
}
