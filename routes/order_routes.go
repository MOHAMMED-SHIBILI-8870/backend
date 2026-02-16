package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine) {
	order := r.Group("/order")
	order.Use(middleware.AuthMiddleware())
	{
		order.POST("/", controllers.PlaceOrder)
		order.GET("/", controllers.GetUserOrders)
		order.GET("/:id", controllers.GetOrder)
	}
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	admin.GET("/orders", controllers.GetAllOrders)
	admin.PUT("/orders/:id/status", controllers.UpdateOrderStatus)
}
