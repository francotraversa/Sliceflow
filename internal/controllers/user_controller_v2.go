package controller

import (
	"net/http"
	"strconv"

	middleware "github.com/francotraversa/Sliceflow/internal/middlewares"
	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userUseCase domain.UserUseCase
}

func NewUserController(userUseCase domain.UserUseCase) *UserController {
	return &UserController{userUseCase: userUseCase}
}

func (c *UserController) CreateUser(ctx echo.Context) error {
	var user types.UserCreateCreds
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, types.Error{Error: err.Error()})
	}
	if err := c.userUseCase.CreateUser(&user, claims.CompanyId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}
	return ctx.JSON(http.StatusOK, user)
}

func (c *UserController) GetUsers(ctx echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, types.Error{Error: err.Error()})
	}
	users, err := c.userUseCase.GetUsers(claims.CompanyId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}
	return ctx.JSON(http.StatusOK, users)
}

func (c *UserController) DeleteUser(ctx echo.Context) error {
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, types.Error{Error: err.Error()})
	}
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	if err := c.userUseCase.DeleteUser(uint(id), claims.CompanyId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}
	return ctx.JSON(http.StatusOK, types.Response{Message: "User deleted successfully"})
}

func (c *UserController) UpdateUser(ctx echo.Context) error {
	var user types.UserUpdateCreds
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, types.Error{Error: err.Error()})
	}
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}
	if err := c.userUseCase.UpdateUser(uint(id), &user, claims.CompanyId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.Error{Error: err.Error()})
	}
	return ctx.JSON(http.StatusOK, user)
}
