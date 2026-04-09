package services

import (
	"fmt"
	"strings"

	userStorage "github.com/francotraversa/Sliceflow/internal/infra/database/user_utils"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllUserUserUseCase(requesterRole string, filterRole string, role string, companyId int) (*[]types.User, error) {
	if requesterRole == "owner" {
		usersDB := userStorage.FindAllUsers(uint(companyId))
		return &usersDB, nil
	}

	filterRole = strings.ToLower(strings.TrimSpace(filterRole))
	if filterRole != "" && filterRole != "admin" && filterRole != "user" {
		return nil, fmt.Errorf("invalid filter role: use 'admin' or 'user'")
	}

	usersDB := userStorage.FindUsersByRole(filterRole, uint(companyId))

	var response []types.User
	for _, u := range usersDB {
		response = append(response, types.User{
			IdUser:    u.IdUser,
			Username:  u.Username,
			Role:      u.Role,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			IdCompany: u.IdCompany,
		})
	}

	return &response, nil
}
