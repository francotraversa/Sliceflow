package machineutils

import (
	"errors"
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

func getDB() *gorm.DB {
	return storage.DatabaseInstance{}.Instance()
}

func GetMachinebyID(id int, companyID uint) (*types.Machine, error) {
	db := getDB()
	var machine types.Machine

	if err := db.Where("id_company = ?", companyID).First(&machine, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Machine doesn't exists")
		}
		return nil, err
	}

	return &machine, nil
}

func GetMachine(dto types.CreateMachineDTO) (*types.Machine, error) {
	db := getDB()
	var machine types.Machine

	// Usamos First
	if err := db.Where("name = ?", dto.Name).First(&machine).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error database lookup for machine %s: %w", dto.Name, err)
	}

	// 3. Si llegamos acá, la máquina existe.
	return &machine, nil
}

func GetMachinesFiltered(filter types.MachineFilter, companyID uint) (*[]types.Machine, error) {
	db := getDB()
	var machines []types.Machine

	query := db.Model(&types.Machine{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Where("id_company = ?", companyID).Find(&machines).Error; err != nil {
		return nil, fmt.Errorf("error database lookup for machines: %w", err)
	}

	return &machines, nil
}
