package repository

import (
	"errors"
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
	// TODO: Implement if needed
	return nil, nil
}

func (r *authRepository) CheckUser(userCreds types.UserLoginCreds) (bool, error) {
	var user types.User
	err := r.db.Where("username = ?", strings.ToLower(strings.TrimSpace(userCreds.Username))).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // usuario no existe → credenciales inválidas
		}
		return false, err
	}
	return true, nil
}

func (r *authRepository) CheckPassword(userCreds types.UserLoginCreds) (bool, error) {
	var user types.User
	err := r.db.Where("username = ?", strings.ToLower(strings.TrimSpace(userCreds.Username))).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if err := utils.CheckPassword(user.Password, userCreds.Password); err != nil {
		return false, nil // contraseña incorrecta → no revelar el tipo de error
	}
	return true, nil
}

func (r *authRepository) GetUser(username string) (*types.User, error) {
	var user types.User
	err := r.db.Where("username = ?", strings.ToLower(strings.TrimSpace(username))).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
