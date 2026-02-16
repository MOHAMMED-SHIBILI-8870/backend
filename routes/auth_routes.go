package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/verify-otp", controllers.VerifyOTPController)
		auth.POST("/login", controllers.Login)
		auth.POST("/forget-pass", controllers.ForgetPassword)
		auth.POST("/reset-pass", controllers.ResetPassword)
		auth.POST("/resent-otp", controllers.ResendOtpHandler)
		auth.POST("/logout",middleware.AdminAuthMiddleware(), controllers.Logout)
	}

}
