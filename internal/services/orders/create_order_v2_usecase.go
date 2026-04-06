package services

import (
	"errors"
	"fmt"
	"time"

	servicesWeb "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *OrderService) CreateOrder(dtoOrder types.CreateOrderDTO, companyID uint) error {
	if dtoOrder.ID == nil {
		return fmt.Errorf("order id is required")
	}
	exists, err := s.repo.GetByID(*dtoOrder.ID, companyID)
	if err != nil {
		return fmt.Errorf("error getting order: %w or does not exist", err)
	}
	if exists {
		return fmt.Errorf("order %d already exists", *dtoOrder.ID)
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
			err := s.machineService.UpdateMachine(uint(*itemDTO.MachineID), updmachine, companyID)

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
		IdCompany:        companyID,
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
	if err := s.repo.Create(&newOrder); err != nil {
		return fmt.Errorf("could not save order and items: %w", err)
	}
	servicesWeb.InvalidateCache("orders:list:*")
	servicesWeb.PublishEvent("dashboard_CreateOrder", `{"type": "ORDER_CREATED", "message": "ORDER CREATED"}`)
	return nil
}
