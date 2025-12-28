package controller

import (
	"net/http"

	services "github.com/francotraversa/Sliceflow/internal/services/authenticator"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// LoginHandler godoc
// @Summary      Iniciar sesión
// @Description  Autentica al usuario y devuelve un token JWT
// @Tags         Authenticator
// @Accept       json
// @Produce      json
// @Param        credentials  body      types.UserLoginCreds  true  "Credenciales"
// @Success      200          {object}  types.TokenResponse
// @Router       /hornero/auth/login [post]
func LoginHandler(c echo.Context) error {
	var userCread types.UserLoginCreds
	if err := c.Bind(&userCread); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid JSON")
	}
	token, err := services.AuthUseCase(userCread)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	return c.JSON(http.StatusOK, token)

}
