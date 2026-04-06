package domain

import "github.com/francotraversa/Sliceflow/internal/types"

type MaterialRepository interface {
	Create(material *types.Material) error
	Update(id uint, material types.UpdateMaterialDTO, companyID uint) error
	Delete(id uint, companyID uint) error
	GetByID(id uint, companyID uint) (*types.Material, error)
	GetMaterials(filter types.MaterialFilter, companyID uint) ([]types.Material, error)
}

// El contrato que tiene que firmar la Mesera (UseCase/Service)
type MaterialUseCase interface {
	CreateMaterial(dto types.CreateMaterialDTO, companyID uint) error
	UpdateMaterial(id uint, dto types.UpdateMaterialDTO, companyID uint) error
	DeleteMaterial(id uint, companyID uint) error
	GetMaterialByID(id uint, companyID uint) (*types.Material, error)
	GetMaterials(filter types.MaterialFilter, companyID uint) ([]types.Material, error)
}
