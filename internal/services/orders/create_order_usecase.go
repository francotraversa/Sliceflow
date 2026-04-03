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
	if dtoOrder.ID == nil {
		return fmt.Errorf("order id is required")
	}
	order, err := ordersutils.CheckOrder(dtoOrder)
	if err != nil {
		return fmt.Errorf("error checking order: %w", err)
	}
	if order != nil {
		return fmt.Errorf("The Order %d already exists", *dtoOrder.ID)
	}
	initialStatus := "pending"
	var itemsDB []types.OrderItem
	totalPiecesCalculated := 0
	totalPriceCalculated := 0.0
	for _, itemDTO := range dtoOrder.Items {
		totalPiecesCalculated += itemDTO.Quantity
		totalPriceCalculated += itemDTO.Price
		itemsDB = append(itemsDB, types.OrderItem{
			StlName:    itemDTO.StlName,
			Quantity:   itemDTO.Quantity,
			DonePieces: 0,
			MaterialID: itemDTO.MaterialID,
			MachineID:  itemDTO.MachineID,
			Price:      &itemDTO.Price,
			Weight:     itemDTO.Weight,
			Time:       itemDTO.Time,
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
		IdOrder:          *dtoOrder.ID,
		IdCompany:        idCompany,
		ClientName:       dtoOrder.ClientName,
		Items:            itemsDB,               // Items list built in the loop
		TotalPieces:      totalPiecesCalculated, // Sum of all item quantities
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
