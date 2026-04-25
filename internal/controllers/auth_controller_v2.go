package controller

import (
	"net/http"

	"github.com/francotraversa/Sliceflow/internal/services/domain"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authUseCase domain.AuthUseCase
}

func NewAuthController(uc domain.AuthUseCase) *AuthController {
	return &AuthController{authUseCase: uc}
}

// Login godoc
// @Summary      Login
// @Description  Authenticates the user and returns a JWT token
// @Tags         Authenticator
// @Accept       json
// @Produce      json
// @Param        credentials  body      types.UserLoginCreds  true  "Credentials"
// @Success      200          {object}  types.TokenResponse
// @Router       /hornero/authed/auth/login [post]
func (ac *AuthController) LoginHandler(c echo.Context) error {
	var userCreds types.UserLoginCreds
	if err := c.Bind(&userCreds); err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}

	token, err := ac.authUseCase.Login(userCreds)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, token)
}
