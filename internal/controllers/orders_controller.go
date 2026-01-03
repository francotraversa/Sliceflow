package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/francotraversa/Sliceflow/internal/auth"
	services "github.com/francotraversa/Sliceflow/internal/services/orders"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// CreateOrderHandler godoc
// @Summary      Crear Orden de Trabajo
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        request body   types.CreateOrderDTO  true  "Formulario Orden"
// @Router       /hornero/authed/orders/order [post]
func CreateOrderHandler(c echo.Context) error {
	var dto types.CreateOrderDTO
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	// Validaciones básicas manuales si querés
	if dto.ClientName == "" || dto.TotalPieces <= 0 || dto.ID == nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Insufficient fields"})
	}

	if err := services.CreateOrderUseCase(dto); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("The Order %d has been created", *dto.ID)})
}

// GetOrdersHandler godoc
// @Summary      Listar Órdenes Activas
// @Tags         Orders
// @Produce      json
// @Router       /hornero/authed/orders/list [get]
func GetOrdersHandler(c echo.Context) error {
	var filter types.OrderFilter

	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Filters don't work"})
	}
	orders, err := services.GetAllOrdersUseCase(filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, orders)
}

// UpdateOrderHandler godoc
// @Summary      Actualizar Orden de Trabajo
// @Description  Permite editar detalles, asignar máquina o actualizar progreso (piezas hechas).
// @Tags         Orders
// @Param        id      path    int                   true  "ID de la Orden"
// @Param        request body    types.UpdateOrderDTO  true  "Datos Nuevos"
// @Router      /hornero/authed/orders/updord/{id} [put]
func UpdateOrderHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "ID param invalid"})
	}

	var dto types.UpdateOrderDTO
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	if err := services.UpdateOrderUseCase(id, dto); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, types.Response{Message: fmt.Sprintf("The Order %d has been updated", id)})
}

// GetDashboardHandler godoc
// @Summary      Dashboard Principal (Role-Based)
// @Description  Muestra métricas y órdenes. Si es admin ve revenue, si no, ve $0.
// @Security     BearerAuth
// @Tags         Production
// @Router       /hornero/authed/orders/dashboard [get]
func GetPrincipalDashboardHandler(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*auth.JwtCustomClaims)

	data, err := services.GetDashboardDataUseCase(claims.Role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, data)
}
