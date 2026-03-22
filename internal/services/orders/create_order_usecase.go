package services

import (
	"errors"
	"fmt"
	"time"

	ordersutils "github.com/francotraversa/Sliceflow/internal/infra/database/orders_utils"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	servicesWeb "github.com/francotraversa/Sliceflow/internal/services/common"
	services "github.com/francotraversa/Sliceflow/internal/services/machine"

	"github.com/francotraversa/Sliceflow/internal/types"
)

func CreateOrderUseCase(dtoOrder types.CreateOrderDTO, idCompany uint) error {
	if dtoOrder.ID != nil {
		order, _ := ordersutils.CheckOrder(dtoOrder)
		if order != nil {
			return fmt.Errorf("The Order %d already exists", *dtoOrder.ID)
		}
	}
	initialStatus := "pending"
	var itemsDB []types.OrderItem
	totalPiecesCalculated := 0
	totalPriceCalculated := 0.0
	for _, itemDTO := range dtoOrder.Items {
		totalPiecesCalculated += itemDTO.Quantity
		totalPriceCalculated += *itemDTO.Price
		itemsDB = append(itemsDB, types.OrderItem{
			StlName:    itemDTO.StlName,
			Quantity:   itemDTO.Quantity,
			DonePieces: 0,
			MaterialID: itemDTO.MaterialID,
			MachineID:  itemDTO.MachineID,
			Price:      itemDTO.Price,
		})

		if itemDTO.MachineID != nil {
			initialStatus = "queued"
		}

		if itemDTO.MachineID != nil {
			newStatus := "printing"
			updmachine := types.UpdateMachineDTO{
				Status: &newStatus,
			}
			err := services.UpdateMachineUseCase(*itemDTO.MachineID, updmachine, idCompany)

			if err != nil {
				return fmt.Errorf("could not update machine status: %w", err)
			}

		}
	}

	deadlineTime, err := time.Parse("2006-01-02", dtoOrder.Deadline)
	if err != nil {
		return errors.New("Format Date invalid (use YYYY-MM-DD)")
	}

	totalMinutes := (dtoOrder.EstimatedHours * 60) + dtoOrder.EstimatedMinutes

	newOrder := types.ProductionOrder{
		ID:               *dtoOrder.ID,
		IdCompany:        idCompany,
		ClientName:       dtoOrder.ClientName,
		Items:            itemsDB,               // La lista de piezas que armamos en el loop
		TotalPieces:      totalPiecesCalculated, // Usamos la suma de las cantidades
		DonePieces:       0,
		Priority:         dtoOrder.Priority,
		Notes:            dtoOrder.Notes,
		EstimatedMinutes: totalMinutes,
		Deadline:         deadlineTime,
		Status:           initialStatus,
		OperatorID:       dtoOrder.OperatorID,
		TotalPrice:       &totalPriceCalculated,
	}

	if err := db_utils.Create(&newOrder); err != nil {
		return fmt.Errorf("could not save order and items: %w", err)
	}

	servicesWeb.InvalidateCache("orders:list:*")
	servicesWeb.PublishEvent("dashboard_updates", `{"type": "ORDER_CREATED", "message": "NEW MULTI-ITEM ORDER CREATED"}`)

	return nil
}
