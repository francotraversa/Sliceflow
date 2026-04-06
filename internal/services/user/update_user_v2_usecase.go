package services

import (
	"fmt"
	"log/slog"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func (u *UserServices) UpdateUser(id uint, userUpdate *types.UserUpdateCreds, companyID uint) error {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		slog.Error("users: user not found", "id", id, "error", err)
		return fmt.Errorf("user not found: %w", err)
	}
	if user.IdCompany != companyID {
		slog.Warn("users: unauthorized update attempt", "target_id", id, "requester_id", companyID)
		return fmt.Errorf("unauthorized: you do not have permission to update this user")
	}
	if err := u.userRepo.UpdateUser(id, userUpdate, companyID); err != nil {
		slog.Error("users: update failed", "id", id, "error", err)
		return fmt.Errorf("failed to update user: %w", err)
	}
	slog.Info("users: updated", "id", id)
	return nil
}
