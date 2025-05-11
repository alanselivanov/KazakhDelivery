package routes

import (
	"api-gateway/internal/config"
	"api-gateway/internal/controllers"
	"api-gateway/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(middlewares.LoggingMiddleware())

	inventoryCtrl := controllers.NewInventoryController(cfg.Services.Inventory)
	orderCtrl := controllers.NewOrderController(cfg.Services.Order)
	userCtrl := controllers.NewUserController(cfg.Services.User)

	products := router.Group("/products")
	{
		products.POST("", inventoryCtrl.CreateProduct)
		products.GET(":id", inventoryCtrl.GetProduct)
		products.PATCH(":id", inventoryCtrl.UpdateProduct)
		products.DELETE(":id", inventoryCtrl.DeleteProduct)
		products.GET("", inventoryCtrl.ListProducts)
	}

	categories := router.Group("/categories")
	{
		categories.POST("", inventoryCtrl.CreateCategory)
		categories.GET(":id", inventoryCtrl.GetCategory)
		categories.PATCH(":id", inventoryCtrl.UpdateCategory)
		categories.DELETE(":id", inventoryCtrl.DeleteCategory)
		categories.GET("", inventoryCtrl.ListCategories)
	}

	orders := router.Group("/orders")
	orders.Use(middlewares.AuthMiddleware())
	{
		orders.POST("", orderCtrl.CreateOrder)
		orders.GET(":id", orderCtrl.GetOrder)
		orders.PATCH(":id", orderCtrl.UpdateOrder)
		orders.GET("", orderCtrl.ListOrders)
	}

	users := router.Group("/users")
	{
		users.POST("/register", userCtrl.RegisterUser)
		users.POST("/login", userCtrl.AuthenticateUser)
		users.GET("/:id/profile", userCtrl.GetUserProfile)
		users.PATCH("/:id/profile", userCtrl.UpdateUserProfile)
	}

	return router
}
