package services

import (
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllMachinesUseCase(filter types.MachineFilter) (*[]types.Machine, error) {
	db := storage.DatabaseInstance{}.Instance()
	cacheKey := fmt.Sprintf("machine:list:%s:%s", filter.Status, filter.Type)

	var machines []types.Machine

	if services.GetCache(cacheKey, &machines) {
		return &machines, nil
	}

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
	services.SetCache(cacheKey, &machines)
	return &machines, nil
}
