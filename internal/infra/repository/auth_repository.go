package repository

import (
	"strings"

	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/francotraversa/Sliceflow/internal/utils"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Login(userCreds types.UserLoginCreds) (*types.TokenResponse, error) {
	// TODO: Implement the repository logic referencing r.db
	return nil, nil
}

func (r *authRepository) CheckUser(userCreds types.UserLoginCreds) (bool, error) {
	var user types.User

	query := r.db.Where("username = ?", strings.ToLower(strings.TrimSpace(userCreds.Username)))
	if err := query.Scan(&user).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *authRepository) CheckPassword(userCreds types.UserLoginCreds) (bool, error) {
	var user types.User

	query := r.db.Where("username = ?", strings.ToLower(strings.TrimSpace(userCreds.Username)))
	if err := query.Scan(&user).Error; err != nil {
		return false, err
	}
	err := utils.CheckPassword(user.Password, userCreds.Password)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *authRepository) GetUser(username string) (*types.User, error) {
	var user types.User

	query := r.db.Where("username = ?", strings.ToLower(strings.TrimSpace(username)))
	if err := query.Scan(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
