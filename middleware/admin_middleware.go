package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Admin_middleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		role,exist:=c.Get("role")
		if !exist || role != "admin"{
			c.JSON(http.StatusForbidden,gin.H{
				"error":"Admin can only access !! : Access denied",
			})
			c.Abort()
			return 
		}
		c.Next()
	}
}