package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func ViewRoutes(r *gin.Engine){
	
	r.GET("/login",controllers.ShowLoginPage)
	view:=r.Group("/view")

	

	view.Use(middleware.AuthMiddleware())
	{
		view.GET("/dashboard",controllers.ShowDashboard)
	}
}