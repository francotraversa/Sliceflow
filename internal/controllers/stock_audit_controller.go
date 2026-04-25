package controller

import (
	"log/slog"
	"net/http"
	"strconv"

	middleware "github.com/francotraversa/Sliceflow/internal/middlewares"
	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

type StockAuditController struct {
	auditService domain.StockAuditService
}

func NewStockAuditController(auditService domain.StockAuditService) *StockAuditController {
	return &StockAuditController{auditService: auditService}
}

func (c *StockAuditController) CreateMovement(ctx echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		slog.Warn("audit: failed to extract JWT claims", "error", err)
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	var req types.CreateMovementRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	req.UserID = claims.UserId

	if err := c.auditService.CreateMovement(req, claims.CompanyId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, types.Response{Message: "Movement created successfully"})
}

func (c *StockAuditController) GetMovementByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid movement ID")
	}

	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	movement, err := c.auditService.GetMovementByID(uint(id), claims.CompanyId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if movement == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Movement not found")
	}

	return ctx.JSON(http.StatusOK, movement)
}

func (c *StockAuditController) GetAllMovements(ctx echo.Context) error {
	var filter types.HistoryFilter
	if err := ctx.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	movements, err := c.auditService.GetAllMovements(filter, claims.CompanyId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, movements)
}
