package controller

import (
	"net/http"

	middleware "github.com/francotraversa/Sliceflow/internal/middlewares"
	services "github.com/francotraversa/Sliceflow/internal/services/metrics"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// GetMetrics godoc
// @Summary      Get metrics
// @Description  Returns metrics for the authenticated user
// @Tags         Metrics
// @Produce      json
// @Security BearerAuth
// @Success      200   {object}  types.MetricsResponse
// @Failure      400   {object}  types.Error
// @Router       /hornero/authed/metrics [get]
func GetMetricsHandler(c echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	response, err := services.GetMetricsUseCase(claims.CompanyId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}
