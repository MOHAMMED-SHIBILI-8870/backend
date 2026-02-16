package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func WishlistRouts(r *gin.Engine) {
	wishlist := r.Group("/wishlist")
	wishlist.Use(middleware.AuthMiddleware())

	{
		wishlist.POST("/", controllers.AddToWishlist)
		wishlist.GET("/", controllers.GetWishlist)
		wishlist.DELETE("/:product_id", controllers.RemoveFromWishlist)
	}
}
