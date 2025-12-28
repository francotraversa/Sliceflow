package main

import (
	"os"

	_ "github.com/francotraversa/Sliceflow/docs"
	"github.com/francotraversa/Sliceflow/internal/auth"
	storage "github.com/francotraversa/Sliceflow/internal/database"
	enviroment "github.com/francotraversa/Sliceflow/internal/enviroment"
	"github.com/francotraversa/Sliceflow/internal/routers"
	"github.com/francotraversa/Sliceflow/internal/swagger"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomClaims struct {
	jwt.RegisteredClaims
}

// -----------------SWAGGER-----------------
// @title           API de Sliceflow
// @version         1.0
// @description     Documentación de mi API con Echo.
// @host            localhost:8181
// @BasePath        /

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                Escribí "Bearer " seguido de tu token JWT. Ejemplo: "Bearer eyJhbG..."
func main() {
	enviroment.LoadEnviroment("dev")
	storage.DatabaseInstance{}.NewDataBase()
	e := echo.New()
	swagger.RegisterSwagger(e)

	e.Use(middleware.Recover())

	jwtCfg := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
		SigningMethod: "HS256",
	}

	routers.RegisterRouters(e, jwtCfg)

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))

}
