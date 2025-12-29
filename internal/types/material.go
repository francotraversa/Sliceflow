package types

import (
	"time"

	"gorm.io/gorm"
)

type MaterialFilter struct {
	Name  string `query:"name"` // ?name=PLA
	Type  string `query:"type"` // ?type=Filamento
	Brand string `json:"brand"`
}

type CreateMaterialDTO struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Brand       string `json:"brand"`
}
type UpdateMaterialDTO struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Brand       string `json:"brand"`
}

type Material struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"type:varchar(100);not null;unique" json:"name"` // Ej: "PLA Negro"
	Type        string `gorm:"type:varchar(50);not null" json:"type"`         // Ej: "Filamento", "Resina"
	Description string `gorm:"type:text" json:"description"`                  // Ej: "Marca Grilon, temperatura 200°C"
	Brand       string `gorm:"type:text" json:"brand"`                        // Ej: "Marca Grilon, temperatura 200°C"
}
