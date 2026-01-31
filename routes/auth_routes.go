package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine){
	auth := r.Group("/auth")
	{
		auth.POST("/register",controllers.Register)
		auth.POST("/verify-otp",controllers.VerifyOTPController)
	}
}