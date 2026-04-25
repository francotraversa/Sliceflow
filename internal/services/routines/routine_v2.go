package services

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	servicesWeb "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func (s *RoutineService) CheckAndSetPriorities() error {
	now := time.Now()

	orders, err := s.repo.GetActiveOrders()
	if err != nil {
		return fmt.Errorf("failed to fetch orders for routine: %w", err)
	}

	var updatedOrders []types.ProductionOrder

	for _, order := range orders {
		if order.Deadline.IsZero() {
			continue
		}

		duration := order.Deadline.Sub(now)
		hoursLeft := duration.Hours()
		hoursStalled := now.Sub(order.UpdatedAt).Hours()

		madeChanges := false

		// 1. Crossover < 0 hrs
		if hoursLeft < 0 {
			if order.Priority != "P1" {
				order.Priority = "P1"
				madeChanges = true
			}
			if !strings.Contains(order.Notes, "[OVERDUE]") {
				order.Notes = strings.TrimSpace(fmt.Sprintf("[OVERDUE] %s", order.Notes))
				madeChanges = true
			}
			if hoursLeft < -24 && order.Status != "late" {
				order.Status = "late"
				madeChanges = true
			}

		} else if hoursLeft < 24 {
			if order.Priority != "P1" {
				order.Priority = "P1"
				madeChanges = true
			}
		} else if hoursLeft < 48 {
			if order.Priority == "P3" || order.Priority == "" {
				order.Priority = "P2"
				madeChanges = true
			}
		}

		// 4. Revisar si están estancadas
		if order.Status == "in-progress" && hoursStalled > 24 {
			if !strings.Contains(order.Notes, "[STALLED?]") {
				order.Notes = strings.TrimSpace(fmt.Sprintf("%s [STALLED?]", order.Notes))
				madeChanges = true
			}
		}

		if madeChanges {
			updatedOrders = append(updatedOrders, order)
		}
	}

	if len(updatedOrders) > 0 {
		if err := s.repo.BulkUpdateOrders(updatedOrders); err != nil {
			return fmt.Errorf("failed to save updated orders: %w", err)
		}

		slog.Info("routine: updated priorities for orders", "count", len(updatedOrders))
		servicesWeb.InvalidateCache("orders:list:*")
		servicesWeb.PublishEvent("dashboard_updates", `{"type": "REFRESH_ORDERS", "message": "New Order List"}`)
	}

	return nil
}
