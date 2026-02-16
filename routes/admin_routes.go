package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {

	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		admin.GET("/users", controllers.GetAllUsersHandler)
		admin.GET("/users/:id", controllers.GetUserByIDHandler)
		admin.PUT("/users/:id", controllers.UpdateUserHandler)
		admin.DELETE("/users/:id", controllers.DeleteUserHandler)
		admin.POST("/users/:id/block", controllers.BlockUserHandler)
		admin.POST("/users/:id/unblock", controllers.UnblockUserHandler)
	}
}