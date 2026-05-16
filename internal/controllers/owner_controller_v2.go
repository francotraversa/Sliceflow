package controller

import (
	"net/http"

	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/labstack/echo/v4"
)

type OwnerController struct {
	ownerService domain.OwnerUseCase
}

func NewOwnerController(ownerService domain.OwnerUseCase) *OwnerController {
	return &OwnerController{ownerService: ownerService}
}

func (c *OwnerController) GetAllUsersHandler(ctx echo.Context) error {

	users, err := c.ownerService.GetAllUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, users)
}
