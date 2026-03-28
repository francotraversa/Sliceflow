package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	middleware "github.com/francotraversa/Sliceflow/internal/middlewares"
	services "github.com/francotraversa/Sliceflow/internal/services/orders"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// CreateOrderHandler godoc
// @Summary      Create Production Order
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        request body   types.CreateOrderDTO  true  "Order Form"
// @Router       /hornero/authed/orders/order [post]
func CreateOrderHandler(c echo.Context) error {
	var dto types.CreateOrderDTO
	if err := c.Bind(&dto); err != nil {
		slog.Warn("orders: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("orders: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	if err := services.CreateOrderUseCase(dto, claims.CompanyId); err != nil {
		slog.Error("orders: creation failed", "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("orders: created", "order_id", *dto.ID, "company_id", claims.CompanyId)
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("The Order %d has been created", *dto.ID)})
}

// GetOrdersHandler godoc
// @Summary      List Active Orders
// @Tags         Orders
// @Produce      json
// @Router       /hornero/authed/orders/list [get]
func GetOrdersHandler(c echo.Context) error {
	var filter types.OrderFilter

	if err := c.Bind(&filter); err != nil {
		slog.Warn("orders: invalid filter params", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Filters invalid"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("orders: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	orders, err := services.GetAllOrdersUseCase(filter, claims.CompanyId)
	if err != nil {
		slog.Error("orders: failed to list", "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("orders: listed", "count", len(*orders), "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, orders)
}

// UpdateOrderHandler godoc
// @Summary      Update Production Order
// @Description  Allows editing details, assigning machines or updating progress (done pieces).
// @Tags         Orders
// @Param        id      path    int                   true  "Order ID"
// @Param        request body    types.UpdateOrderDTO  true  "Updated data"
// @Router      /hornero/authed/orders/updord/{id} [put]
func UpdateOrderHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("orders: invalid ID param", "param", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "ID param invalid"})
	}

	var dto types.UpdateOrderDTO
	if err := c.Bind(&dto); err != nil {
		slog.Warn("orders: invalid request body", "order_id", id, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("orders: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	if err := services.UpdateOrderUseCase(id, dto, claims.CompanyId); err != nil {
		slog.Error("orders: update failed", "order_id", id, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("orders: updated", "order_id", id, "company_id", claims.CompanyId)
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("The Order %d has been updated", id)})
}

// GetDashboardHandler godoc
// @Summary      Dashboard Principal (Role-Based)
// @Description  Shows metrics and orders. Admins see revenue, others see $0.
// @Security     BearerAuth
// @Tags         Production
// @Router       /hornero/authed/orders/dashboard [get]
func GetPrincipalDashboardHandler(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*types.JwtCustomClaims)

	data, err := services.GetDashboardDataUseCase(claims.Role, claims.CompanyId)
	if err != nil {
		slog.Error("orders: dashboard failed", "company_id", claims.CompanyId, "role", claims.Role, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("orders: dashboard retrieved", "company_id", claims.CompanyId, "role", claims.Role)
	return c.JSON(http.StatusOK, data)
}

// DeleteOrderHandler godoc
// @Summary      Delete Production Order
// @Description  Deletes an order by its ID. Admin only.
// @Tags         Orders
// @Param        id  path    int true "Order ID"
// @Security     BearerAuth
// @Router       /hornero/authed/orders/delord/{id} [delete]
func DeleteOrderHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("orders: invalid ID param", "param", c.Param("id"), "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "ID param invalid"})
	}
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		slog.Warn("orders: failed to extract JWT claims", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	if err := services.DeleteOrderUseCase(id, claims.CompanyId); err != nil {
		slog.Error("orders: deletion failed", "order_id", id, "company_id", claims.CompanyId, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("orders: deleted", "order_id", id, "company_id", claims.CompanyId)
	return c.JSON(http.StatusOK, types.Response{Message: fmt.Sprintf("The Order %d has been deleted", id)})
}
