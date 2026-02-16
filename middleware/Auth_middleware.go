package middleware

import (
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		access_token, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			c.Abort()
			return
		}

		userID, role, err := utils.ValidateJwt(access_token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		if role == "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "admins are not allowed",
			})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Set("role", role)

		c.Next()
	}
}
