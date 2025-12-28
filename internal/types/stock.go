package types

import (
	"time"

	"gorm.io/gorm"
)

type StockItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	SKU         string         `gorm:"uniqueIndex;not null" json:"sku"`
	Name        string         `gorm:"not null" json:"name"`
	Quantity    int            `gorm:"default:0" json:"quantity"`
	MinQty      float64        `gorm:"default:5" json:"min_qty"`
	Description string         `gorm:"null" json:"description"`
	Status      string         `gorm:"size:16;default:'active'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type StockLocation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type StockLevel struct {
	ID         uint          `gorm:"primaryKey" json:"id"`
	ItemID     uint          `gorm:"uniqueIndex:idx_item_location" json:"item_id"`
	LocationID uint          `gorm:"uniqueIndex:idx_item_location" json:"location_id"`
	Qty        float64       `gorm:"default:0" json:"qty"`
	Item       StockItem     `gorm:"foreignKey:ItemID" json:"-"`
	Location   StockLocation `gorm:"foreignKey:LocationID" json:"-"`
}

type StockMovement struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ItemID     uint      `gorm:"index" json:"item_id"`
	LocationID uint      `gorm:"index" json:"location_id"`
	Type       string    `json:"type"` // IN, OUT, ADJUST
	QtyDelta   float64   `json:"qty_delta"`
	QtyBefore  float64   `json:"qty_before"`
	QtyAfter   float64   `json:"qty_after"`
	Reason     string    `json:"reason"`
	CreatedBy  uint      `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type ProductCreateRequest struct {
	SKU         string `json:"sku" example:"7791234567890" extensions:"x-order=1"`
	Name        string `json:"name" example:"PLA PRINT A LOT 1KG AZUL"`
	Description string `json:"description" example:"Importación Diciembre"`
	Quantity    int    `json:"Quantity" example:"10"`
}

type ProductUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}
