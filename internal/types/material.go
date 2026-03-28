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

	Name        string  `gorm:"type:varchar(100);not null;unique" json:"name"` // e.g. "PLA Black"
	Type        string  `gorm:"type:varchar(50);not null" json:"type"`         // e.g. "Filament", "Resin"
	Description string  `gorm:"type:text" json:"description"`                  // e.g. "Grilon brand, 200°C temp"
	Brand       string  `gorm:"type:text" json:"brand"`                        // e.g. "Grilon"
	IdCompany   uint     `gorm:"not null" json:"id_company"`
	Company     *Company `gorm:"foreignKey:IdCompany;references:IdCompany" json:"company,omitempty"`
}
