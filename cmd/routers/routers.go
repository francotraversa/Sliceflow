package routers

import (
	controller "github.com/francotraversa/Sliceflow/internal/controllers"
	storage "github.com/francotraversa/Sliceflow/internal/infra/database"
	"github.com/francotraversa/Sliceflow/internal/middlewares"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RegisterRouters(e *echo.Echo, jwtCfg echojwt.Config) {

	db := storage.DBInstance.Instance()

	AppControllers := BuildDependencies(db)

	//---------------------PUBLIC---------------------------------
	e.GET("/health", controller.RegisterHealth)

	api := e.Group("/hornero")
	auth := api.Group("/auth")
	auth.POST("/login", AppControllers.Auth.LoginHandler)
	//----------------------WEBSOCKET (JWT via query param) ----
	api.GET("/ws/dashboard", controller.WebSocketHandler)
	//----------------------LOGED------------------------------
	protected := api.Group("/authed")
	protected.Use(echojwt.WithConfig(jwtCfg))
	//----------------------USER---------------------------------
	user := protected.Group("/user")
	user.PATCH("/updmyuser", controller.UpdateUserHandler)
	user.DELETE("/delmyuser", controller.DeleteUserHandler)
	//---------------------STOCK PRODUCTS-------------------------
	stock := protected.Group("/stock")
	stock.POST("/addprod", controller.CreateProductHandler)
	stock.GET("/list", controller.GetProductsHandler)
	stock.PATCH("/updprod/:sku", controller.UpdateByIdProductHandler)
	stock.DELETE("/delprod/:sku", controller.DeleteIdProductHandler)
	//----------------------STOCK MOVEMENT-----------------------
	movement := stock.Group("/movement")
	movement.POST("/addmov", controller.CreateMovementHandler)
	movement.GET("/historic", controller.GetStockHistoryHandler)
	movement.GET("/dashboard", controller.GetDashboardHandler)
	//----------------------MATERIAL---------------------------
	production := protected.Group("/materials")
	production.POST("/addmat", AppControllers.Material.CreateMaterialHandler)
	production.PATCH("/updmat/:id", AppControllers.Material.UpdateMaterialHandler)
	production.DELETE("/delmat/:id", AppControllers.Material.DeleteMaterialHandler)
	production.GET("/list", AppControllers.Material.GetMaterialsHandler)
	//----------------------MACHINE-------------------------------
	machine := protected.Group("/machine")
	machine.POST("/addmac", AppControllers.Machine.CreateMachineHandler)
	machine.PATCH("/updmac/:id", AppControllers.Machine.UpdateMachineHandler)
	machine.GET("/list", AppControllers.Machine.GetMachinesHandler)
	machine.DELETE("/delmac/:id", AppControllers.Machine.DeleteMachineHandler)
	//----------------------ORDERS------------------------------
	orders := protected.Group("/orders")
	orders.POST("/order", controller.CreateOrderHandler)
	orders.GET("/list", controller.GetOrdersHandler)
	orders.PATCH("/updord/:id", controller.UpdateOrderHandler)
	orders.GET("/dashboard", controller.GetPrincipalDashboardHandler)
	orders.DELETE("/delord/:id", controller.DeleteOrderHandler)
	orders.GET("/metrics", AppControllers.Metrics.GetMetricsHandler)
	//----------------------PRIVATE && ROLE----------------------
	admin := protected.Group("/admin")
	admin.Use(middlewares.RequireRole("admin"))
	admin.POST("/newuser", controller.CreateUserHandler)
	admin.GET("/alluser", controller.GetAllUserHandler)
	admin.PATCH("/edituser/:id", controller.UpdateUserHandler)
	admin.PATCH("/enableuser", controller.EnableUserHandler)
	admin.DELETE("/deleteuser/:id", controller.DeleteUserHandler)
	//----------------------OWNER------------------------------
	owner := protected.Group("/owner")
	owner.Use(middlewares.RequireRole("owner"))
	owner.POST("/newcompany", controller.CreateCompanyHandler)
	owner.POST("/newadmin", controller.CreateAdminHandler)
	owner.DELETE("/deladmin/:id", controller.DeleteAdminHandler)
	owner.GET("/alladmin", controller.GetAllAdminHandler)
	owner.GET("/allcompany", controller.GetAllCompanyHandler)
	owner.DELETE("/deletecompany/:id", controller.DeleteCompanyHandler)
}
