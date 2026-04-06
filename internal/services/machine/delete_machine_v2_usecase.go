package services

import (
	"fmt"
)

func (s *machineService) DeleteMachine(id uint, companyID uint) error {
	machine, err := s.repo.GetByID(id, companyID)
	if err != nil {
		return err
	}
	if machine.IdCompany != companyID {
		return fmt.Errorf("You don't have permission to delete this machine")
	}
	return s.repo.Delete(id, companyID)
}
