package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func ViewRoutes(r *gin.Engine) {

	r.GET("/login", controllers.ShowLoginPage)
	view := r.Group("/view")

	view.Use(middleware.AdminAuthMiddleware())
	{
		view.GET("/dashboard", controllers.ShowDashboard)

		view.GET("/users", controllers.ShowUsersPage)
		view.GET("/products", controllers.ShowProductsPage)
		view.GET("/orders", controllers.ShowOrdersPage)
		//---------USER EDTITE
		view.GET("/users/edit/:id", controllers.ShowEditUserPage)
		//---------- PRODUCT CREATE & UPDATE
		view.GET("/products/create", controllers.ShowCreateProductPage)
		view.GET("/products/edit/:id", controllers.ShowEditProductPage)

		// ---------- ADMIN PROFILE ----------
		view.GET("/profile", controllers.ShowAdminProfilePage)
		view.GET("/profile/edit", controllers.ShowEditAdminProfilePage)
		view.POST("/profile/update", controllers.UpdateAdminProfile)
	}
}
