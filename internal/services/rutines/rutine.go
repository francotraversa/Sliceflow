package services

import (
	"fmt"
	"strings"
	"time"

	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	db_utils "github.com/francotraversa/Sliceflow/internal/infra/database/utils"
	services "github.com/francotraversa/Sliceflow/internal/services/common"
	"github.com/francotraversa/Sliceflow/internal/types"
)

func CheckAndSetPriorities() {
	db := storage.DatabaseInstance{}.Instance()
	now := time.Now()

	var orders []types.ProductionOrder

	if err := db.Find(&orders).Error; err != nil {
		fmt.Println("Error getting Orders in Rutines:", err)
		return
	}

	updatedCount := 0

	for _, order := range orders {
		if order.Deadline.IsZero() {
			continue
		}

		duration := order.Deadline.Sub(now)
		hoursLeft := duration.Hours()

		hoursStalled := now.Sub(order.UpdatedAt).Hours()

		madeChanges := false

		if hoursLeft < 0 {
			if order.Priority != "P1" {
				order.Priority = "P1"
				madeChanges = true
			}
			if !strings.Contains(order.Notes, "[VENCIDA]") {
				order.Notes = strings.TrimSpace(fmt.Sprintf("[VENCIDA] %s", order.Notes))
				madeChanges = true
			}

			// B. URGENTE (< 24 hs)
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

		if order.Status == "in-progress" && hoursStalled > 24 {
			if !strings.Contains(order.Notes, "[ESTANCADA?]") {
				order.Notes = strings.TrimSpace(fmt.Sprintf("%s [ESTANCADA?]", order.Notes))
				madeChanges = true
			}
		}

		if madeChanges {
			if err := db_utils.Save(&order); err != nil {
				fmt.Printf("Error guardando orden %d: %v\n", order.ID, err)
				continue
			}
			updatedCount++
		}
	}

	if updatedCount > 0 {
		services.InvalidateCache("orders:list:*")
		services.PublishEvent("dashboard_updates", `{"type": "REFRESH_ORDERS", "message": "New Order List"}`)
	}
}
