package services

import (
	"github.com/francotraversa/Sliceflow/internal/services/domain"
)

type machineService struct {
	repo domain.MachineRepository
}

func NewMachineService(repo domain.MachineRepository) domain.MachineUseCase {
	return &machineService{repo: repo}
}
