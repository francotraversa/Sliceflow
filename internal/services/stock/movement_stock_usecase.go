package services

import (
	"errors"
	"fmt"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddStockMovementUseCase(req types.CreateMovementRequest, companyID uint) error {
	db := storage.DatabaseInstance{}.Instance()

	// START TRANSACTION
	// Everything inside this block is atomic. Either both changes are saved, or neither.
	return db.Transaction(func(tx *gorm.DB) error {
		var item types.StockItem

		// 1. BUSCAR Y BLOQUEAR (Locking)
		// clause.Locking{Strength: "UPDATE"} le dice a la DB:
		// "Lock this row for writing until the transaction is complete".
		// Evita condiciones de carrera (Race Conditions).
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, "sku = ? AND id_company = ?", req.SKU, companyID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("product with SKU '%s' not found", req.SKU)
			}
			return err
		}
		if item.Id == 0 {
			return fmt.Errorf("product with SKU '%s' has invalid ID (schema migration may be pending)", req.SKU)
		}

		// 2. PREPARE AUDIT DATA
		qtyBefore := item.Quantity
		var qtyDelta int

		// 3. MOVEMENT LOGIC
		switch req.Type {
		case "IN", "RETURN":
			qtyDelta = req.Quantity

		case "OUT", "INTERNAL_USE":
			qtyDelta = -req.Quantity // Lo volvemos negativo para la suma

			// CRITICAL validation: Don't sell what we don't have
			if item.Quantity < req.Quantity {
				return fmt.Errorf("stock insuficiente. Tienes %d, intentas sacar %d", item.Quantity, req.Quantity)
			}
		case "ADJUSTMENT":
			qtyDelta = req.Quantity - item.Quantity

		default:
			return errors.New("invalid movement type (use IN or OUT)")
		}

		qtyAfter := qtyBefore + qtyDelta

		// 4. ACTUALIZAR EL OBJETO PRODUCTO (En memoria)
		item.Quantity = qtyAfter

		// 5. CREAR EL OBJETO MOVIMIENTO
		movement := types.StockMovement{
			StockItemID: item.Id,  // FK by DB Id (resolved from SKU + companyID)
			StockSKU:    item.SKU, // kept for display
			Type:        req.Type,
			QtyDelta:    qtyDelta,
			QtyBefore:   qtyBefore,
			QtyAfter:    qtyAfter,
			Description: req.Description,
			IdCompany:   companyID,
			Reason:      req.Reason,
			CreatedBy:   req.UserID,
			LocationID:  req.LocationID,
		}

		// 6. GUARDAR CAMBIOS EN LA DB (Usando 'tx', no 'db')

		if err := tx.Create(&movement).Error; err != nil {
			return err // Dispara Rollback
		}

		// B. Guardamos el producto actualizado
		// Usamos Select("Quantity") para asegurarnos de solo tocar el stock y no pisar otros datos
		if err := tx.Model(&item).Select("Quantity").Updates(&item).Error; err != nil {
			return err // Dispara Rollback
		}
		services.InvalidateCache(fmt.Sprintf("stock:list:%d", companyID))
		services.InvalidateCache(fmt.Sprintf("historic:list:*company=%d", companyID))
		services.InvalidateCache("dashboard:*")
		services.PublishEvent("dashboard_updates", `{"type": "STOCK_MOVEMENT", "message": "STOCK MOVEMENT CREATED"}`)
		return nil
	})
}
