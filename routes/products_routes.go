package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	
	{
		admin.POST("/createproduct",middleware.AdminAuthMiddleware(), controllers.CreateProduct)
		admin.PUT("/updateproduct/:id",middleware.AdminAuthMiddleware(), controllers.UpdateProduct)
		admin.DELETE("/deleteproduct/:id", middleware.AdminAuthMiddleware(),controllers.DeleteProduct)
	}
	public:=r.Group("/products")
	{
		public.GET("/:id",controllers.GetProductByID)
		public.GET("/",controllers.GetAllProducts)
	}
}
