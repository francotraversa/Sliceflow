package services

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/francotraversa/Sliceflow/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserServices) CreateUser(user *types.UserCreateCreds, companyID uint) error {
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

	existingUser, err := s.userRepo.GetUserByUsername(strings.ToLower(user.Username))

	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if existingUser != nil {
		return fmt.Errorf("user with username %s already exists", strings.ToLower(user.Username))
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
		Username:  strings.ToLower(user.Username),
		Password:  string(hash),
		Role:      user.Role,
		IdCompany: *user.IdCompany,
	}

	if err := s.userRepo.CreateUser(&u); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	slog.Info("users: created", "username", user.Username)
	return nil
}
