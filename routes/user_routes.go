package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserProfileRoutes(r *gin.Engine){
	user:=r.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/profile",controllers.GetProfile)
		user.PUT("/profile",controllers.UpdateProfile)
	}
}