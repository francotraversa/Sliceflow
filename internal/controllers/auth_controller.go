package controller

import (
	"log/slog"
	"net/http"

	services "github.com/francotraversa/Sliceflow/internal/services/authenticator"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/labstack/echo/v4"
)

// LoginHandler godoc
// @Summary      Login
// @Description  Authenticates the user and returns a JWT token
// @Tags         Authenticator
// @Accept       json
// @Produce      json
// @Param        credentials  body      types.UserLoginCreds  true  "Credentials"
// @Success      200          {object}  types.TokenResponse
// @Router       /hornero/auth/login [post]
func LoginHandler(c echo.Context) error {
	var userCread types.UserLoginCreds
	if err := c.Bind(&userCread); err != nil {
		slog.Warn("auth: invalid request body", "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: "Invalid Json"})
	}
	token, err := services.AuthUseCase(userCread)
	if err != nil {
		slog.Warn("auth: login failed", "username", userCread.Username, "error", err)
		return c.JSON(http.StatusBadRequest, types.Error{Error: err.Error()})
	}

	slog.Info("auth: login successful", "username", userCread.Username)
	return c.JSON(http.StatusOK, token)
}
