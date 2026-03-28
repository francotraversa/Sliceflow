package controller

import (
	"net/http"

	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// RegisterHealth godoc
// @Summary      Health check
// @Description  Returns OK if the service is up
// @Tags         Health
// @Produce      json
// @Success      200  {string}  string "OK"
// @Router       /health [get]
func RegisterHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, types.Response{Message: "ok"})
}
