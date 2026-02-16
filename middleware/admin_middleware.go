package middleware

import (
	"backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// authHeader := c.GetHeader("Authorization")
		access_token,err:=c.Cookie("access_token")
		if err != nil{
			log.Fatal(err)
		}
		log.Println(access_token +"recieved")

		// if authHeader == "" {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		// 	return
		// }

		// tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		userId, role, err := utils.ValidateJwt(access_token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return

		}

		if role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user",
			})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Set("role", role)

		c.Next()
	}
}