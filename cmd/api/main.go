package main

import (
	"log"
	"os"
	"time"

	_ "github.com/francotraversa/Sliceflow/docs"
	"github.com/francotraversa/Sliceflow/internal/auth"
	enviroment "github.com/francotraversa/Sliceflow/internal/enviroment"
	redis "github.com/francotraversa/Sliceflow/internal/infra/cache"
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	userStorage "github.com/francotraversa/Sliceflow/internal/infra/database/user_utils"
	"github.com/francotraversa/Sliceflow/internal/routers"
	services "github.com/francotraversa/Sliceflow/internal/services/rutines"
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
// @host            localhost:1000
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name 					   Authorization
// @description                Escribí "Bearer " seguido de tu token JWT. Ejemplo: "Bearer eyJhbG..."
func main() {
	enviroment.LoadEnviroment("dev")
	storage.DatabaseInstance{}.NewDataBase()
	if err := userStorage.EnsureHardcodedUser(); err != nil {
		log.Fatalf("Error creando usuario hardcodeado: %v", err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redis.InitRedis(redisHost, "6379", "")

	e := echo.New()
	e.Use(middleware.Recover())

	jwtCfg := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
		SigningMethod: "HS256",
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

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))

}
