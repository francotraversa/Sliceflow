package services

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *machineService) UpdateMachine(id uint, dto types.UpdateMachineDTO, companyID uint) error {
	machine, err := s.repo.GetByID(id, companyID)
	if err != nil {
		return err
	}

	if dto.Name != nil {
		machine.Name = *dto.Name
	}

	if dto.Type != nil {
		machine.Type = *dto.Type
	}

	if dto.Status != nil {
		machine.Status = *dto.Status
	}

	if err := s.repo.Update(id, companyID, dto); err != nil {
		return fmt.Errorf("The Machine was not updated")
	}
	return nil
}
