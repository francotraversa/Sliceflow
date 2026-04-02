package services

import (
	"fmt"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	userStorage "github.com/francotraversa/Sliceflow/internal/infra/database/user_utils"
	"github.com/francotraversa/Sliceflow/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserUseCase(user types.UserCreateCreds, companyID uint) error {
	db := storage.DatabaseInstance{}.Instance()

	if user.Username == "" || user.Password == "" {
		return fmt.Errorf("Username or password field is empty")
	}
	if len(user.Password) < 6 {
		return fmt.Errorf("short password (min 6)")
	}

	if user.Role == "" {
		user.Role = "user"
	} else {
		normalized := strings.ToLower(strings.TrimSpace(user.Role))
		if normalized == "user" {
			user.Role = normalized
		} else {
			return fmt.Errorf("invalid user role")
		}
	}

	usercheck := userStorage.FindUserByUsername(user.Username)

	if usercheck != nil {
		return fmt.Errorf("The user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error generating password hash")
	}
	if user.IdCompany == nil {
		user.IdCompany = &companyID
	}

	if user.Role == "admin" {
		if *user.IdCompany != companyID {
			return fmt.Errorf("You can only create users for your own company")
		}
	}

	u := types.User{
		Username:  user.Username,
		Password:  string(hash),
		Role:      user.Role,
		IdCompany: *user.IdCompany,
	}
	if err := db.Create(&u).Error; err != nil {
		return fmt.Errorf("Error creating user")
	}
	return nil

}
