package main

import (
	"backend/config"
	"backend/routes"

	"backend/models"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{},&models.OTP{})

	r := gin.Default()
	
	routes.AuthRoutes(r)

	r.Run(":8080")
}
