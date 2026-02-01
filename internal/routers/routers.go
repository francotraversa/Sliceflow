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
	api.GET("/ws/dashboard", controller.WebSocketHandler)
	auth := api.Group("/auth")
	auth.POST("/login", controller.LoginHandler)
	//----------------------PRIVATE------------------------------
	protected := api.Group("/authed")
	protected.Use(echojwt.WithConfig(jwtCfg))
	//----------------------USER---------------------------------
	user := protected.Group("/user")
	user.PATCH("/updmyuser", controller.UpdateUserHandler)
	user.DELETE("/delmyuser", controller.DeleteUserHandler)
	//---------------------STOCK PRODUCTS-------------------------
	stock := protected.Group("/stock")
	stock.POST("/addprod", controller.CreateProductHandler)
	stock.GET("/list", controller.GetAllProductsHandler)
	stock.GET("/list/:name", controller.GetProductByNameHandler)
	stock.GET("/list/:sku", controller.GetIdProductHandler)
	stock.PATCH("/updprod/:sku", controller.UpdateByIdProductHandler)
	stock.DELETE("/delprod/:sku", controller.DeleteIdProductHandler)
	//----------------------STOCK MOVEMENT-----------------------
	movement := stock.Group("/movement")
	movement.POST("/addmov", controller.CreateMovementHandler)
	movement.GET("/historic", controller.GetStockHistoryHandler)
	movement.GET("/dashboard", controller.GetDashboardHandler)
	//----------------------MATERIAL---------------------------
	production := protected.Group("/materials")
	production.POST("/addmat", controller.CreateMaterialHandler)
	production.PATCH("/updmat/:id", controller.UpdateMaterialHandler)
	production.DELETE("/delmat/:id", controller.DeleteMaterialHandler)
	production.GET("/list", controller.GetMaterialsHandler)
	//----------------------MACHINE-------------------------------
	machine := protected.Group("/machine")
	machine.POST("/addmac", controller.CreateMachineHandler)
	machine.PATCH("/updmac/:id", controller.UpdateMachineHandler)
	machine.GET("/list", controller.GetMachinesHandler)
	machine.DELETE("/delmac/:id", controller.DeleteMachineHandler)
	//----------------------ORDERS------------------------------
	orders := protected.Group("/orders")
	orders.POST("/order", controller.CreateOrderHandler)
	orders.GET("/list", controller.GetOrdersHandler)
	orders.PATCH("/updord/:id", controller.UpdateOrderHandler)
	orders.GET("/dashboard", controller.GetPrincipalDashboardHandler)
	//----------------------PRIVATE && ROLE----------------------
	admin := protected.Group("/admin")
	admin.Use(middlewares.RequireRole("admin"))
	admin.POST("/newuser", controller.CreateUserHandler)
	admin.GET("/alluser", controller.GetAllUserHandler)
	admin.PATCH("/edituser/:id", controller.UpdateUserHandler)
	admin.PATCH("/enableuser", controller.EnableUserHandler)
	admin.DELETE("/deleteuser/:id", controller.DeleteUserHandler)

}
