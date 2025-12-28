package types

import (
	"time"

	"gorm.io/gorm"
)

type StockItem struct {
	SKU         string         `gorm:"primaryKey;size:50;not null" json:"sku"`
	Name        string         `gorm:"not null;index" json:"name"`
	Quantity    int            `gorm:"default:0" json:"quantity"`
	Price       float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	MinQty      float64        `gorm:"default:5" json:"min_qty"`
	Description string         `gorm:"null" json:"description"`
	Status      string         `gorm:"size:16;default:'active';index" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type StockMovement struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	StockSKU   string    `gorm:"size:50;index;not null" json:"stock_sku"`
	StockItem  StockItem `gorm:"foreignKey:StockSKU;references:SKU" json:"-"`
	LocationID uint      `gorm:"index" json:"location_id"`
	Type       string    `gorm:"size:20;not null" json:"type"` // IN, OUT, ADJUST
	QtyDelta   int       `gorm:"not null" json:"qty_delta"`    // Cuánto cambió (+10, -5)
	QtyBefore  int       `gorm:"not null" json:"qty_before"`   // Cuánto había (100)
	QtyAfter   int       `gorm:"not null" json:"qty_after"`    // Cuánto quedó (110)
	Reason     string    `gorm:"size:255" json:"reason"`
	CreatedBy  uint      `gorm:"index" json:"created_by"`
	User       User      `gorm:"foreignKey:CreatedBy" json:"-"` // Opcional: relación con User
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
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

type ProductCreateRequest struct {
	SKU         string  `json:"sku" example:"7791234567890" extensions:"x-order=1"`
	Name        string  `json:"name" example:"PLA PRINT A LOT 1KG AZUL"`
	Description string  `json:"description" example:"Importación Diciembre"`
	Quantity    int     `json:"Quantity" example:"10"`
	Price       float64 `json:"price" example:"98.5"`
}

type ProductUpdateRequest struct {
	SKU         string  `json:"sku" example:"7791234567890" extensions:"x-order=1"`
	Name        string  `json:"name" example:"PLA PRINT A LOT 1KG AZUL"`
	Description string  `json:"description" example:"Importación Diciembre"`
	Quantity    int     `json:"quantity" example:"10"`
	Price       float64 `json:"price" example:"98.5"`
}

type CreateMovementRequest struct {
	SKU        string `json:"sku" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required,gt=0"`     // Siempre positivo, el Tipo define si suma o resta
	Type       string `json:"type" validate:"required,oneof=IN OUT"` // Solo permitimos IN o OUT
	Reason     string `json:"reason"`
	UserID     uint   `json:"user_id"`     // Viene del token/contexto
	LocationID uint   `json:"location_id"` // Opcional por ahora
}

type HistoryFilter struct {
	SKU       string `query:"sku"`
	StartDate string `query:"start_date"` // Formato esperado: YYYY-MM-DD
	EndDate   string `query:"end_date"`   // Formato esperado: YYYY-MM-DD
}
