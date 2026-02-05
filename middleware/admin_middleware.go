package middleware

import (
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		userId,role, err := utils.ValidateJwt(tokenStr)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		if role != "admin" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("userId", userId)
		c.Set("role", role)

		c.Next()
	}
}
