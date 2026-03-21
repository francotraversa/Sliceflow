package types

import (
	"time"

	"gorm.io/gorm"
)

type OrderFilter struct {
	ID           uint   `query:"id"`            // ?id=123
	Status       string `query:"status"`        // ?status=pending
	SortPriority bool   `query:"sort_priority"` // ?sort_priority=true
	FromDate     string `query:"from_date"`     // ?from_date=2024-01-01
	ToDate       string `query:"to_date"`       // ?to_date=2024-01-31
}

type ProductionOrder struct {
	ID               uint           `gorm:"primaryKey;autoIncrement:false" json:"id"` // autoIncrement:false es clave para ID manual
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	ClientName       string         `gorm:"type:varchar(150);not null" json:"client_name"`
	Items            []OrderItem    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"items"`
	Priority         string         `gorm:"type:varchar(10);default:'P3'" json:"priority"`
	Notes            string         `gorm:"type:text" json:"notes"`
	TotalPieces      int            `gorm:"not null" json:"total_pieces"`
	EstimatedMinutes int            `json:"estimated_minutes"`
	Deadline         time.Time      `gorm:"type:timestamp;not null" json:"deadline"`
	FinishTime       *time.Time     `json:"finish_time,omitempty"`
	Status           string         `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	DonePieces       int            `gorm:"default:0" json:"done_pieces"`
	OperatorID       int            `gorm:"not null" json:"operator_id"`
	Price            *float64       `gorm:"type:decimal(12,2)" json:"price,omitempty"`
}

type OrderItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `gorm:"index" json:"order_id"`
	StlName    string    `gorm:"type:varchar(150);not null" json:"product_name"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	DonePieces int       `gorm:"default:0" json:"done_pieces"`
	MaterialID *int      `gorm:"index" json:"material_id"`
	MachineID  *int      `gorm:"index" json:"machine_id"`
	Material   *Material `gorm:"foreignKey:MaterialID" json:"material"`
	Machine    *Machine  `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
}

type CreateOrderDTO struct {
	ID               *uint                `gorm:"primaryKey" json:"id"`
	ClientName       string               `json:"client_name"`
	Items            []CreateOrderItemDTO `json:"items"` // Lista de piezas
	MaterialID       *int                 `json:"material_id"`
	Priority         string               `json:"priority"`
	Notes            string               `json:"notes"`
	EstimatedHours   int                  `json:"estimated_hours"`
	EstimatedMinutes int                  `json:"estimated_minutes"`
	Deadline         string               `json:"deadline"` // "YYYY-MM-DD"
	OperatorID       int                  `json:"operator_id"`
	Price            *float64             `json:"price,omitempty"`
}
type CreateOrderItemDTO struct {
	ID          uint   `json:"id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	DonePieces  int    `json:"done_pieces"`
	MachineID   *int   `json:"machine_id"`
	MaterialID  *int   `json:"material_id"`
}

type UpdateOrderDTO struct {
	ID               *uint                 `json:"id"`
	ClientName       *string               `json:"client_name"`
	ProductDetails   *string               `json:"product_details"`
	TotalPieces      *int                  `json:"total_pieces"`
	DonePieces       *int                  `json:"done_pieces"`
	Priority         *string               `json:"priority"` // O int, segun como lo tengas
	Notes            *string               `json:"notes"`
	Status           *string               `json:"status"`
	Price            *float64              `json:"price"`
	EstimatedHours   *int                  `json:"estimated_hours"`
	EstimatedMinutes *int                  `json:"estimated_minutes"`
	Deadline         *string               `json:"deadline"`
	Items            *[]CreateOrderItemDTO `json:"items" `

	OperatorID *int `json:"operator_id"`
	MaterialID *int `json:"material_id"`
	MachineID  *int `json:"machine_id"`
}
