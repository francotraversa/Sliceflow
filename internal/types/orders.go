package types

import (
	"time"

	"gorm.io/gorm"
)

type OrderFilter struct {
	Status       string `query:"status"`        // ?status=pending
	SortPriority bool   `query:"sort_priority"` // ?sort_priority=true
}

// ProductionOrder: La orden de trabajo
type ProductionOrder struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Datos del Cliente
	ClientName     string `gorm:"type:varchar(150);not null" json:"client_name"`
	ProductDetails string `gorm:"type:text" json:"product_details"` // Texto largo por si detallan mucho

	// Especificaciones
	TotalPieces int       `gorm:"not null" json:"total_pieces"`
	MaterialID  int       `gorm:"not null" json:"material_id"`                   // La Foreign Key
	Material    *Material `gorm:"foreignKey:MaterialID" json:"material"`         // La Relación (Objeto)
	Priority    string    `gorm:"type:varchar(10);default:'P3'" json:"priority"` // Default prioridad baja
	Notes       string    `gorm:"type:text" json:"notes"`

	// Tiempos
	EstimatedMinutes int        `json:"estimated_minutes"`
	Deadline         time.Time  `json:"deadline"`
	StartDate        *time.Time `json:"start_date"` // Puntero porque puede ser nulo al inicio

	// Estado y Progreso
	Status     string `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	DonePieces int    `gorm:"default:0" json:"done_pieces"`

	// Relaciones (Foreign Keys)
	OperatorID int `gorm:"not null" json:"operator_id"`

	MachineID *int    `gorm:"index" json:"machine_id"`
	Machine   Machine `gorm:"foreignKey:MachineID" json:"machine,omitempty"`

	Price float64 `gorm:"type:decimal(12,2)" json:"price"`
}
type CreateOrderDTO struct {
	ClientName     string `json:"client_name"`
	ProductDetails string `json:"product_details"`
	TotalPieces    int    `json:"total_pieces"`

	MaterialID int `json:"material_id"` // <--- El ID del combo box

	Priority string `json:"priority"` // "P1", "P2", "P3"
	Notes    string `json:"notes"`

	EstimatedHours   int `json:"estimated_hours"`   // input horas
	EstimatedMinutes int `json:"estimated_minutes"` // input minutos

	Deadline string `json:"deadline"` // String "YYYY-MM-DD"

	OperatorID int     `json:"operator_id"` // El ID del usuario logueado o seleccionado
	MachineID  *int    `json:"machine_id"`  // Puede ser null (sin asignar)
	Price      float64 `json:"price"`
}

type UpdateOrderDTO struct {
	ClientName       *string  `json:"client_name"`
	ProductDetails   *string  `json:"product_details"`
	TotalPieces      *int     `json:"total_pieces"`
	DonePieces       *int     `json:"done_pieces"`
	Priority         *string  `json:"priority"` // O int, segun como lo tengas
	Notes            *string  `json:"notes"`
	Status           *string  `json:"status"`
	Price            *float64 `json:"price"`
	EstimatedHours   *int     `json:"estimated_hours"`
	EstimatedMinutes *int     `json:"estimated_minutes"`
	Deadline         *string  `json:"deadline"`

	// Relaciones (Foreign Keys)
	OperatorID *int `json:"operator_id"`
	MaterialID *int `json:"material_id"`
	MachineID  *int `json:"machine_id"`
}
