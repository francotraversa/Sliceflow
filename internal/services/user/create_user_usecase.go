package services

import (
	"fmt"
	"strings"

	"github.com/francotraversa/Sliceflow/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserUseCase(userCreds types.UserCreateCreds) error {
	// 1. Obtenemos la DB de la instancia global que ya está conectada
	db := database.DBInstance.DB

	if userCreds.Username == "" || userCreds.Password == "" {
		return fmt.Errorf("username or password field is empty")
	}
	if len(userCreds.Password) < 6 {
		return fmt.Errorf("short password (min 6)")
	}

	if userCreds.Role == "" {
		userCreds.Role = "user"
	} else {
		userCreds.Role = strings.ToLower(strings.TrimSpace(userCreds.Role))
		if userCreds.Role != "user" && userCreds.Role != "admin" {
			return fmt.Errorf("invalid user role")
		}
	}

	// 2. Usamos el nuevo paquete 'user' (alias userDB o el que prefieras)
	// Pasamos db como primer argumento para evitar el import cycle
	usercheck := user_utils.FindUserByUsername(db, userCreds.Username)

	if usercheck != nil {
		return fmt.Errorf("the user already exists")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(userCreds.Password), bcrypt.DefaultCost)

	u := types.User{
		Username: userCreds.Username,
		Password: string(hash),
		Role:     userCreds.Role,
	}

	if err := db.Create(&u).Error; err != nil {
		return fmt.Errorf("error creando usuario")
	}
	return nil
}
