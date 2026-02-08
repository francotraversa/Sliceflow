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

func AddStockMovementUseCase(req types.CreateMovementRequest) error {
	db := storage.DatabaseInstance{}.Instance()

	// INICIO DE TRANSACCIÓN
	// Todo lo que pase acá adentro es atómico. O se guardan los dos cambios, o ninguno.
	return db.Transaction(func(tx *gorm.DB) error {
		var item types.StockItem

		// 1. BUSCAR Y BLOQUEAR (Locking)
		// clause.Locking{Strength: "UPDATE"} le dice a la DB:
		// "Bloquea esta fila para escritura hasta que termine la transacción".
		// Evita condiciones de carrera (Race Conditions).
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, "sku = ?", req.SKU).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("producto no encontrado")
			}
			return err
		}

		// 2. PREPARAR DATOS DE AUDITORÍA
		qtyBefore := item.Quantity
		var qtyDelta int

		// 3. LÓGICA DE MOVIMIENTO
		switch req.Type {
		case "IN", "RETURN":
			qtyDelta = req.Quantity

		case "OUT", "INTERNAL_USE":
			qtyDelta = -req.Quantity // Lo volvemos negativo para la suma

			// Validación CRÍTICA: No vender lo que no tenemos
			if item.Quantity < req.Quantity {
				return fmt.Errorf("stock insuficiente. Tienes %d, intentas sacar %d", item.Quantity, req.Quantity)
			}
		case "ADJUSTMENT":
			qtyDelta = req.Quantity - item.Quantity

		default:
			return errors.New("tipo de movimiento inválido (Use IN o OUT)")
		}

		qtyAfter := qtyBefore + qtyDelta

		// 4. ACTUALIZAR EL OBJETO PRODUCTO (En memoria)
		item.Quantity = qtyAfter

		// 5. CREAR EL OBJETO MOVIMIENTO
		movement := types.StockMovement{
			StockSKU:    item.SKU, // Relación por string
			Type:        req.Type,
			QtyDelta:    qtyDelta,  // ej: +5 o -5
			QtyBefore:   qtyBefore, // ej: 100
			QtyAfter:    qtyAfter,  // ej: 105 o 95
			Description: req.Description,

			Reason:     req.Reason,
			CreatedBy:  req.UserID,
			LocationID: req.LocationID,
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
		services.InvalidateCache("stock:list:*")
		services.InvalidateCache("historic:list:*")
		services.InvalidateCache("dashboard:*")
		services.PublishEvent("dashboard_updates", `{"type": "STOCK_MOVEMENT", "message": "STOCK MOVEMENT CREATED"}`)
		return nil
	})
}
