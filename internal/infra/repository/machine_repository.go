package repository

import (
	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

type machineRepository struct {
	db *gorm.DB
}

func NewMachineRepository(db *gorm.DB) domain.MachineRepository {
	return &machineRepository{db: db}
}
func (r *machineRepository) Create(machine *types.Machine) error {
	return r.db.Create(machine).Error
}

func (r *machineRepository) Update(id uint, companyID uint, dto types.UpdateMachineDTO) error {
	return r.db.Where("id = ? AND id_company = ?", id, companyID).Updates(dto).Error
}

func (r *machineRepository) Delete(id uint, companyID uint) error {
	return r.db.Where("id = ? AND id_company = ?", id, companyID).Delete(&types.Machine{}).Error
}

func (r *machineRepository) GetByID(id uint, companyID uint) (*types.Machine, error) {
	var machine types.Machine
	if err := r.db.Where("id = ? AND id_company = ?", id, companyID).First(&machine).Error; err != nil {
		return nil, err
	}
	return &machine, nil
}

func (r *machineRepository) GetMachines(filter types.MachineFilter, companyID uint) ([]types.Machine, error) {
	var machines []types.Machine

	query := r.db.Where("id_company = ?", companyID)
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if err := query.Find(&machines).Error; err != nil {
		return nil, err
	}
	return machines, nil
}
