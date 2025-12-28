package services

import (
	"fmt"
	"strings"

	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllUserUserUseCase(requesterRole string, filterRole string) ([]types.User, error) {
	if requesterRole != "admin" {
		return nil, fmt.Errorf("permission denied: only admins can access the user list")
	}

	filterRole = strings.ToLower(strings.TrimSpace(filterRole))
	if filterRole != "" && filterRole != "admin" && filterRole != "user" {
		return nil, fmt.Errorf("invalid filter role: use 'admin' or 'user'")
	}

	usersDB := storage.FindUsersByRole(filterRole)

	var response []types.User
	for _, u := range usersDB {
		response = append(response, types.User{
			ID:       u.ID,
			Username: u.Username,
			Role:     u.Role,
			Status:   u.Status,
		})
	}

	return response, nil
}
