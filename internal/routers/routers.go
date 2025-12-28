package routers

import (
	controller "github.com/francotraversa/Sliceflow/internal/controllers"
	"github.com/francotraversa/Sliceflow/internal/middlewares"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RegisterRouters(e *echo.Echo, jwtCfg echojwt.Config) {
	//---------------------PUBLIC---------------------------------
	e.GET("/health", controller.RegisterHealth)

	api := e.Group("/hornero")
	auth := api.Group("/auth")
	auth.POST("/login", controller.LoginHandler)

	//----------------------PRIVATE------------------------------
	protected := api.Group("/loged")
	protected.Use(echojwt.WithConfig(jwtCfg))
	protected.PATCH("/updmyuser", controller.UpdateUserHandler)
	protected.DELETE("/delmyuser", controller.DeleteUserHandler)

	//----------------------PRIVATE && ROLE----------------------
	admin := protected.Group("/admin")
	admin.Use(middlewares.RequireRole("admin"))
	admin.POST("/newuser", controller.CreateUserHandler)
	admin.GET("/alluser", controller.GetAllUserHandler)
	admin.PATCH("/edituser/:id", controller.UpdateUserHandler)
	admin.DELETE("/deleteuser/:id", controller.DeleteUserHandler)

}
