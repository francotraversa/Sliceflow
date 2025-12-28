package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterHealth godoc
// @Summary      Chequear estado
// @Description  Devuelve OK si el servicio está arriba
// @Tags         Health
// @Produce      json
// @Success      200  {string}  string "OK"
// @Router       /health [get]
func RegisterHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
