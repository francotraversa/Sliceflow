package repository

import (
	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
)

type materialRepository struct {
	db *gorm.DB
}

func NewMaterialRepository(db *gorm.DB) domain.MaterialRepository {
	return &materialRepository{db: db}
}

func (r *materialRepository) Create(material *types.Material) error {
	return r.db.Create(material).Error
}

func (r *materialRepository) Update(id uint, material types.UpdateMaterialDTO, companyID uint) error {
	return r.db.Where("id = ? AND id_company = ?", id, companyID).Updates(&material).Error
}

func (r *materialRepository) Delete(id uint, companyID uint) error {
	return r.db.Where("id = ? AND id_company = ?", id, companyID).Delete(&types.Material{}).Error
}

func (r *materialRepository) GetByID(id uint, companyID uint) (*types.Material, error) {
	var material types.Material
	err := r.db.Where("id = ? AND id_company = ?", id, companyID).First(&material).Error
	return &material, err
}

func (r *materialRepository) GetMaterials(filter types.MaterialFilter, companyID uint) ([]types.Material, error) {
	var materials []types.Material

	query := r.db.Where("id_company = ?", companyID)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.Brand != "" {
		query = query.Where("brand = ?", filter.Brand)
	}

	if err := query.Find(&materials).Error; err != nil {
		return nil, err
	}

	return materials, nil
}
