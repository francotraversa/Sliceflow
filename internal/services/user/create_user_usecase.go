package services

import (
	"fmt"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	userStorage "github.com/francotraversa/Sliceflow/internal/infra/database/user_utils"
	"github.com/francotraversa/Sliceflow/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserUseCase(user types.UserCreateCreds) error {
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
		if user.Role == "user" || user.Role == "admin" {
			user.Role = strings.ToLower(strings.TrimSpace(user.Role))
		} else {
			return fmt.Errorf("invalid user role")
		}
	}

	usercheck := userStorage.FindUserByUsername(user.Username)

	if usercheck != nil {
		return fmt.Errorf("The user already exists")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	u := types.User{
		Username: user.Username,
		Password: string(hash),
		Role:     user.Role,
	}
	if err := db.Create(&u).Error; err != nil {
		return fmt.Errorf("Error creating user")
	}
	return nil

}
