package services

import (
	"fmt"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/francotraversa/Sliceflow/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func UpdateUserUseCase(targetID uint, requesterID uint, requesterRole string, data types.UserUpdateCreds) error {
	db := storage.DatabaseInstance{}.Instance()

	if requesterID != targetID && requesterRole != "admin" {
		return fmt.Errorf("permission denied: you can only update your own profile")
	}

	currentUser := storage.FindUserByUserId(targetID)
	if currentUser == nil {
		return fmt.Errorf("user with ID %d not found", targetID)
	}

	hasChanged := false

	newUsername := strings.ToLower(strings.TrimSpace(data.Username))
	if newUsername != "" && newUsername != currentUser.Username {
		userConflict := storage.FindUserByUsername(newUsername)
		if userConflict != nil {
			return fmt.Errorf("the username '%s' is already taken", newUsername)
		}
		currentUser.Username = newUsername
		hasChanged = true
	}

	if data.Password != "" {
		err := utils.CheckPassword(currentUser.Password, data.Password)
		if err != nil {
			if len(data.Password) < 6 {
				return fmt.Errorf("new password too short (min 6 characters)")
			}
			hash, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
			currentUser.Password = string(hash)
			hasChanged = true
		}
	}

	newRole := strings.ToLower(strings.TrimSpace(data.Role))
	if newRole != "" && newRole != currentUser.Role {
		if requesterRole != "admin" {
			return fmt.Errorf("permission denied: only admins can change user roles")
		}

		if newRole == "user" || newRole == "admin" {
			currentUser.Role = newRole
			hasChanged = true
		} else {
			return fmt.Errorf("invalid role: must be 'admin' or 'user'")
		}
	}

	if !hasChanged {
		fmt.Println("No changes detected for user:", targetID)
		return nil
	}

	if err := db.Save(&currentUser).Error; err != nil {
		return fmt.Errorf("failed to update user in database")
	}

	fmt.Printf("User %d updated successfully by %d\n", targetID, requesterID)
	return nil
}
