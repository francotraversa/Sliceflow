package repository

import (
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

type ownerRepository struct {
	db *gorm.DB
}

func NewOwnerRepository(db *gorm.DB) *ownerRepository {
	return &ownerRepository{db: db}
}

func (r *ownerRepository) GetAllUsers() (*[]types.User, error) {
	var users []types.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}
