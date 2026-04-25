package services

import (
	"fmt"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *machineService) CreateMachine(dto types.CreateMachineDTO, companyID uint) error {
	if dto.Name == "" {
		return fmt.Errorf("Name is required")
	}
	if dto.Type == "" {
		return fmt.Errorf("Type is required")
	}

	newMachine := types.Machine{
		Name:      dto.Name,
		Type:      dto.Type,
		Status:    "idle",
		IdCompany: companyID,
	}
	return s.repo.Create(&newMachine)
}
