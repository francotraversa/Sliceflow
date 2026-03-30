package userStorage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func EnsureHardcodedUser() error {
	db := storage.DatabaseInstance{}.Instance()
	userAdmin := os.Getenv("USERADMIN")
	passAdmin := os.Getenv("PASSWORDADMIN")
	roleAdmin := os.Getenv("ROLEADMIN")
	companyName := os.Getenv("COMPANYADMIN")

	if userAdmin == "" || passAdmin == "" || roleAdmin == "" || companyName == "" {
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

	// 4. Create or find the company
	var company types.Company
	err = db.Where("LOWER(name) = ?", strings.ToLower(companyName)).First(&company).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		company = types.Company{
			Name: companyName,
		}
		if err := db.Create(&company).Error; err != nil {
			log.Printf("[seed] Error creating company '%s': %v", companyName, err)
			return fmt.Errorf("Error creating company")
		}
		log.Printf("[seed] Company '%s' created with ID %d.", companyName, company.IdCompany)
	} else if err != nil {
		return fmt.Errorf("Error looking up company on DB")
	} else {
		log.Printf("[seed] Company '%s' already exists (ID %d).", companyName, company.IdCompany)
	}

	// 5. Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(passAdmin), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[seed] Error hashing password: %v", err)
		return fmt.Errorf("Error hashing password")
	}

	// 6. Create the new admin user with the company
	newUser := types.User{
		Username:  userAdmin,
		Password:  string(hash),
		Role:      roleAdmin,
		IdCompany: company.IdCompany,
	}

	if err := db.Create(&newUser).Error; err != nil {
		log.Printf("[seed] Error creating admin on DB: %v", err)
		return err
	}

	log.Printf("[seed] User '%s' created as admin for company '%s'.", userAdmin, companyName)
	return nil
}
