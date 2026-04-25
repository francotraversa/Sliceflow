package routers

import (
	controller "github.com/francotraversa/Sliceflow/internal/controllers"
	"github.com/francotraversa/Sliceflow/internal/infra/repository"
	authService "github.com/francotraversa/Sliceflow/internal/services/authenticator"
	machineService "github.com/francotraversa/Sliceflow/internal/services/machine"
	materialService "github.com/francotraversa/Sliceflow/internal/services/material"
	metricsService "github.com/francotraversa/Sliceflow/internal/services/metrics"
	orderService "github.com/francotraversa/Sliceflow/internal/services/orders"
	stockService "github.com/francotraversa/Sliceflow/internal/services/stock"
	userService "github.com/francotraversa/Sliceflow/internal/services/user"
	"gorm.io/gorm"
)

type AppControllers struct {
	Auth       *controller.AuthController
	Machine    *controller.MachineController
	Material   *controller.MaterialController
	Metrics    *controller.MetricsController
	Order      *controller.OrderController
	User       *controller.UserController
	Stock      *controller.StockController
	StockAudit *controller.StockAuditController
}

func BuildDependencies(db *gorm.DB) *AppControllers {
	mUseCase := machineService.NewMachineService(repository.NewMachineRepository(db))
	stockRepo := repository.NewStockRepository(db)
	return &AppControllers{
		Auth:       controller.NewAuthController(authService.NewAuthUseCase(repository.NewAuthRepository(db))),
		Machine:    controller.NewMachineController(mUseCase),
		Material:   controller.NewMaterialController(materialService.NewMaterialService(repository.NewMaterialRepository(db))),
		Metrics:    controller.NewMetricsController(metricsService.NewMetricsService(repository.NewMetricsRepository(db))),
		Order:      controller.NewOrderController(orderService.NewOrderService(repository.NewOrderRepository(db), mUseCase)),
		User:       controller.NewUserController(userService.NewUserServices(repository.NewUserRepository(db))),
		Stock:      controller.NewStockController(stockService.NewStockService(stockRepo)),
		StockAudit: controller.NewStockAuditController(stockService.NewStockAuditService(repository.NewStockAuditRepository(db))),
	}
}
