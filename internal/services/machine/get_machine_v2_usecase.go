package services

import (
	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *machineService) GetMachines(filter types.MachineFilter, companyID uint) ([]types.Machine, error) {
	return s.repo.GetMachines(filter, companyID)
}

func (s *machineService) GetMachineByID(id uint, companyID uint) (*types.Machine, error) {
	return s.repo.GetByID(id, companyID)
}
