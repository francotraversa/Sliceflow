package userStorage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/francotraversa/Sliceflow/internal/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func EnsureHardcodedUser(db *gorm.DB) error {
	userAdmin := os.Getenv("USERADMIN")
	passAdmin := os.Getenv("PASSWORDADMIN")
	roleAdmin := os.Getenv("ROLEADMIN")

	if userAdmin == "" || passAdmin == "" {
		log.Println("[seed] WARNING: USERADMIN or PASSWORDADMIN are not definitive. Ommiting seeding...")
		return fmt.Errorf("Error Get Env")
	}

	a := strings.ToLower(userAdmin)

	var existingUser types.User
	err := db.Where("LOWER(username) = ?", a).First(&existingUser).Error

	if err == nil {
		log.Printf("[seed] The user admin '%s' already exist. No changes have been commited.", userAdmin)
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("Error on DB")
	}

	// 4. Hashear la contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(passAdmin), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[seed] Error hashing password: %v", err)
		return fmt.Errorf("Error hashing password")
	}

	// 5. Crear el nuevo usuario admin
	newUser := types.User{
		Username: userAdmin,
		Password: string(hash),
		Role:     roleAdmin,
	}

	if err := db.Create(&newUser).Error; err != nil {
		log.Printf("[seed] Error creating admin on DB: %v", err)
		return err
	}

	log.Printf("[seed] ¡Congratulation! User '%s' has been creadted as admin.", userAdmin)
	return nil
}
