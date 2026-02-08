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
	ID          uint      `gorm:"primaryKey" json:"id"`
	StockSKU    string    `gorm:"size:50;index;not null" json:"stock_sku"`
	StockItem   StockItem `gorm:"foreignKey:StockSKU;references:SKU" json:"-"`
	LocationID  uint      `gorm:"index" json:"location_id"`
	Type        string    `gorm:"size:20;not null" json:"type"` // IN, OUT, ADJUST
	QtyDelta    int       `gorm:"not null" json:"qty_delta"`    // Cuánto cambió (+10, -5)
	QtyBefore   int       `gorm:"not null" json:"qty_before"`   // Cuánto había (100)
	QtyAfter    int       `gorm:"not null" json:"qty_after"`    // Cuánto quedó (110)
	Reason      string    `gorm:"size:255" json:"reason"`
	CreatedBy   uint      `gorm:"index" json:"created_by"`
	User        User      `gorm:"foreignKey:CreatedBy" json:"-"` // Opcional: relación con User
	CreatedAt   time.Time `gorm:"index" json:"created_at"`
	Description string    `gorm:"null" json:"description"`
}

type ProductCreateRequest struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	MinQty      float64 `json:"min_qty"`
}

type ProductUpdateRequest struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	MinQty      float64 `json:"min_qty"`
}

type CreateMovementRequest struct {
	SKU         string `json:"sku" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required,gt=0"`
	Type        string `json:"type" validate:"required,oneof=IN OUT ADJUSTMENT INTERNAL_USE"`
	Reason      string `json:"reason"`
	Description string `json:"description"` // ✅ add this
	UserID      uint   `json:"user_id"`
	LocationID  uint   `json:"location_id"`
}

type HistoryFilter struct {
	SKU       string `query:"sku"`
	StartDate string `query:"start_date"` // Formato esperado: YYYY-MM-DD
	EndDate   string `query:"end_date"`   // Formato esperado: YYYY-MM-DD
	Type      string `query:"type"`
}

type TopProduct struct {
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	TotalSold int    `json:"total_sold"`
}

// DashboardResponse es el objeto gigante que recibe el frontend
type DashboardResponse struct {
	TotalItems      int64        `json:"total_items"`     // Cantidad de productos en catálogo
	TotalValue      float64      `json:"total_value"`     // Cuánta plata hay parada en el depósito
	LowStockCount   int64        `json:"low_stock_count"` // Cuántos productos están en alerta
	MovementsToday  int64        `json:"movements_today"`
	ActiveUsers     int64        `json:"active_users"`
	LowStockItems   []StockItem  `json:"low_stock_items"`   // La lista de esos productos
	TopSellingItems []TopProduct `json:"top_selling_items"` // Los 5 más vendidos
}
