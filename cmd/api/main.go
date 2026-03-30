package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/francotraversa/Sliceflow/docs"
	enviroment "github.com/francotraversa/Sliceflow/internal/environment"
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	userStorage "github.com/francotraversa/Sliceflow/internal/infra/database/user_utils"
	"github.com/francotraversa/Sliceflow/internal/routers"
	services "github.com/francotraversa/Sliceflow/internal/services/routines"
	"github.com/francotraversa/Sliceflow/internal/swagger"
	"github.com/francotraversa/Sliceflow/internal/types"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// -----------------SWAGGER-----------------
// @title           API de Sliceflow
// @version         1.0
// @description     Documentación de mi API con Echo.
// @host            localhost:1000
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name 					   Authorization
// @description                Escribí "Bearer " seguido de tu token JWT. Ejemplo: "Bearer eyJhbG..."
func main() {
	enviroment.LoadEnvironment("dev")
	storage.DatabaseInstance{}.NewDataBase()
	//redis.InitRedis()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{fmt.Sprintf("%s", os.Getenv("FRONTENDHOST"))},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.PATCH},
	}))

	e.Use(middleware.Recover())

	jwtCfg := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(types.JwtCustomClaims)
		},
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
		SigningMethod: "HS256",
		// WebSocket no puede enviar el header Authorization desde el browser:
		// busca el token también en ?token= como fallback
		TokenLookup: "header:Authorization,query:token",
	}
	swagger.RegisterSwagger(e)
	routers.RegisterRouters(e, jwtCfg)

	// ---------------------------------------------------------
	// 🤖 GOROUTINE
	// ---------------------------------------------------------
	go func() {
		time.Sleep(1 * time.Minute)
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			services.CheckAndSetPriorities()
		}
	}()
	if err := userStorage.EnsureHardcodedUser(); err != nil {
		log.Fatalf("Error creando usuario hardcodeado: %v", err)
	}

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))

}
