package services

import (
	"fmt"

	machineutils "github.com/francotraversa/Sliceflow/internal/infra/database/machine_utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func GetAllMachinesUseCase(filter types.MachineFilter, companyID uint) (*[]types.Machine, error) {
	cacheKey := fmt.Sprintf("machine:list:%s:%s", filter.Status, filter.Type)

	var machines []types.Machine

	if services.GetCache(cacheKey, &machines) {
		return &machines, nil
	}

	result, err := machineutils.GetMachinesFiltered(filter, companyID)
	if err == nil {
		machines = *result
	} else {
		return nil, fmt.Errorf("error database lookup for machines: %w", err)
	}

	services.SetCache(cacheKey, &machines)
	return &machines, nil
}
