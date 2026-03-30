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
	Id               uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	IdOrder          uint           `gorm:"not null;uniqueIndex:idx_order_company" json:"id_order"` // business key, unique per company
	IdCompany        uint           `gorm:"not null;uniqueIndex:idx_order_company" json:"id_company"`
	Company          *Company       `gorm:"foreignKey:IdCompany;references:IdCompany" json:"company,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	ClientName       string         `gorm:"type:varchar(150);not null" json:"client_name"`
	Items            []OrderItem    `gorm:"foreignKey:OrderID;references:Id;constraint:OnDelete:CASCADE;" json:"items"`
	Priority         string         `gorm:"type:varchar(10);default:'P3'" json:"priority"`
	Notes            string         `gorm:"type:text" json:"notes"`
	TotalPieces      int            `gorm:"not null" json:"total_pieces"`
	EstimatedMinutes int            `json:"estimated_minutes"`
	Deadline         time.Time      `gorm:"type:timestamp;not null" json:"deadline"`
	FinishTime       *time.Time     `json:"finish_time,omitempty"`
	Status           string         `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	DonePieces       int            `gorm:"default:0" json:"done_pieces"`
	OperatorID       int            `gorm:"not null" json:"operator_id"`
	TotalPrice       *float64       `gorm:"type:decimal(12,2)" json:"total_price,omitempty"`
}

type OrderItem struct {
	ID         uint      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	OrderID    uint      `gorm:"index" json:"order_id"`
	StlName    string    `gorm:"type:varchar(150);not null" json:"product_name"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	DonePieces int       `gorm:"default:0" json:"done_pieces"`
	Price      *float64  `gorm:"type:decimal(12,2)" json:"price,omitempty"`
	MaterialID *int      `gorm:"index" json:"material_id"`
	MachineID  *int      `gorm:"index" json:"machine_id"`
	Material   *Material `gorm:"foreignKey:MaterialID" json:"material"`
	Machine    *Machine  `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
}

type CreateOrderDTO struct {
	ID               *uint                `json:"id"`
	ClientName       string               `json:"client_name"`
	Items            []CreateOrderItemDTO `json:"items"` // List of pieces/parts
	Priority         string               `json:"priority"`
	Notes            string               `json:"notes"`
	EstimatedHours   int                  `json:"estimated_hours"`
	EstimatedMinutes int                  `json:"estimated_minutes"`
	Deadline         string               `json:"deadline"` // "YYYY-MM-DD"
	OperatorID       int                  `json:"operator_id"`
	TotalPrice       *float64             `json:"total_price,omitempty"`
}
type CreateOrderItemDTO struct {
	ID         uint    `json:"id"`
	StlName    string  `json:"stl_name"`
	Quantity   int     `json:"quantity"`
	DonePieces int     `json:"done_pieces"`
	MachineID  *int    `json:"machine_id"`
	MaterialID *int    `json:"material_id"`
	Price      float64 `json:"price"`
}

type UpdateOrderDTO struct {
	ID               *uint                 `json:"id"`
	ClientName       *string               `json:"client_name"`
	ProductDetails   *string               `json:"product_details"`
	TotalPieces      *int                  `json:"total_pieces"`
	DonePieces       *int                  `json:"done_pieces"`
	Priority         *string               `json:"priority"` // Can also be int depending on your setup
	Notes            *string               `json:"notes"`
	Status           *string               `json:"status"`
	TotalPrice       *float64              `json:"total_price"`
	EstimatedHours   *int                  `json:"estimated_hours"`
	EstimatedMinutes *int                  `json:"estimated_minutes"`
	Deadline         *string               `json:"deadline"`
	Items            *[]CreateOrderItemDTO `json:"items" `

	OperatorID *int `json:"operator_id"`
}
