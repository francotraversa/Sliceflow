package types

import (
	"time"

	"gorm.io/gorm"
)

type MachineFilter struct {
	Status string `query:"status"` // ?status=maintenance
	Type   string `query:"type"`   // ?type=FDM
}

type Machine struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name      string   `gorm:"type:varchar(100);not null" json:"name"`        // Required name
	Type      string   `gorm:"type:varchar(50);not null" json:"type"`         // FDM or SLS
	Status    string   `gorm:"type:varchar(20);default:'idle'" json:"status"` // Starts as idle
	IdCompany uint     `gorm:"not null" json:"id_company"`
	Company   *Company `gorm:"foreignKey:IdCompany;references:IdCompany"`
}
type CreateMachineDTO struct {
	Name string `json:"name"` // e.g. "Prusa i3 MK3S+"
	Type string `json:"type"` // e.g. "FDM", "SLS"
}

// DTO for updating a machine (includes Status)
type UpdateMachineDTO struct {
	Name   *string `json:"name"`
	Type   *string `json:"type"`
	Status *string `json:"status"` // "idle", "printing", "maintenance"
}
