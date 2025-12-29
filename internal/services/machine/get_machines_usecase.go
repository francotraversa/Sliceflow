package services

import (
	storage "github.com/francotraversa/Sliceflow/internal/database"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllMachinesUseCase(filter types.MachineFilter) (*[]types.Machine, error) {
	db := storage.DatabaseInstance{}.Instance()
	var machines []types.Machine

	query := db.Model(&types.Machine{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if err := query.Find(&machines).Error; err != nil {
		return nil, err
	}
	return &machines, nil
}
