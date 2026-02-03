package middleware

import (
	"backend/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		header:=c.GetHeader("Authorization")

		if header == ""{
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"Athorization header missing",
			})
			c.Abort()
			return 
		}

		tokenStr:=strings.Replace(header,"","Bearer ",1)
		claims := jwt.MapClaims{}

		token,err:=jwt.ParseWithClaims(tokenStr,claims,func(t *jwt.Token) (any, error) {
			return []byte(config.GetEnv("JWT_ACCESS","access")),nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"Invalid credentials",
			})
			c.Abort()
			return 
		}

		c.Set("user_id",uint(claims["user_id"].(float64)))
		c.Set("role",claims["role"])
		c.Next()
	}
}