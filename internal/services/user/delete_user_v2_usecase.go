package services

import (
	"fmt"
	"log/slog"
)

func (u *UserServices) DeleteUser(id uint, companyID uint) error {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		slog.Error("users: user not found", "id", id, "error", err)
		return fmt.Errorf("user not found: %w", err)
	}
	if user.IdCompany != companyID {
		slog.Warn("users: unauthorized deletion attempt", "target_id", id, "requester_id", companyID)
		return fmt.Errorf("unauthorized: you do not have permission to delete this user")
	}
	if err := u.userRepo.DeleteUser(id, companyID); err != nil {
		slog.Error("users: deletion failed", "id", id, "error", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	slog.Info("users: deleted", "id", id)
	return nil
}
