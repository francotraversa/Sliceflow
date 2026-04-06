package domain

import "github.com/francotraversa/Sliceflow/internal/types"

// El contrato que tiene que firmar el Cocinero (Repository)
type MachineRepository interface {
	Create(machine *types.Machine) error
	Update(id uint, companyID uint, dto types.UpdateMachineDTO) error
	Delete(id uint, companyID uint) error
	GetByID(id uint, companyID uint) (*types.Machine, error)
	GetMachines(filter types.MachineFilter, companyID uint) ([]types.Machine, error)
}

// El contrato que tiene que firmar la Mesera (UseCase/Service)
type MachineUseCase interface {
	CreateMachine(dto types.CreateMachineDTO, companyID uint) error
	UpdateMachine(id uint, dto types.UpdateMachineDTO, companyID uint) error
	DeleteMachine(id uint, companyID uint) error
	GetMachineByID(id uint, companyID uint) (*types.Machine, error)
	GetMachines(filter types.MachineFilter, companyID uint) ([]types.Machine, error)
}
