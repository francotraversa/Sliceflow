package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	middleware "github.com/francotraversa/Sliceflow/internal/middlewares"
	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

type OrderController struct {
	useCase domain.OrderUseCase
}

func NewOrderController(useCase domain.OrderUseCase) *OrderController {
	return &OrderController{useCase: useCase}
}

func (c *OrderController) CreateOrderHandler(ctx echo.Context) error {
	var order types.CreateOrderDTO
	if err := ctx.Bind(&order); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}
	if err := c.useCase.CreateOrder(order, claims.CompanyId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "Order created successfully"})
}

func (c *OrderController) UpdateOrderHandler(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		slog.Warn("orders: invalid ID param", "param", ctx.Param("id"), "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: "ID param invalid"})
	}

	var dto types.UpdateOrderDTO
	if err := ctx.Bind(&dto); err != nil {
		slog.Warn("orders: invalid request body", "order_id", id, "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		slog.Warn("orders: failed to extract JWT claims", "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	if err := c.useCase.UpdateOrder(uint(id), dto, claims.CompanyId); err != nil {
		slog.Error("orders: update failed", "order_id", id, "company_id", claims.CompanyId, "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("orders: updated", "order_id", id, "company_id", claims.CompanyId)
	return ctx.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("The Order %d has been updated", id)})
}

func (c *OrderController) GetOrdersByStatusHandler(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		slog.Warn("orders: invalid ID param", "param", ctx.Param("id"), "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: "ID param invalid"})
	}

	var dto types.OrderFilter
	if err := ctx.Bind(&dto); err != nil {
		slog.Warn("orders: invalid request body", "order_id", id, "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		slog.Warn("orders: failed to extract JWT claims", "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	if orders, err := c.useCase.GetOrdersByStatus(dto, claims.CompanyId); err != nil {
		slog.Error("orders: update failed", "order_id", id, "company_id", claims.CompanyId, "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, orders)
	}

}

func (c *OrderController) DashboardOrdersHandler(ctx echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		slog.Warn("orders: failed to extract JWT claims", "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	if dashboardItems, err := c.useCase.DashboardOrders(claims.Role, claims.CompanyId); err != nil {
		slog.Error("orders: update failed", "company_id", claims.CompanyId, "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, dashboardItems)
	}

}
