package main

import (
	"backend/config"
	"backend/routes"
	"backend/seeder"
	"log"

	"backend/models"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()
	models.Migrate()

	err := seeder.AdminSeeder(config.DB)

	if err != nil{
		log.Fatalf("admin seeder is failed:%v",err)
	}

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/static","./static")
	
	routes.AuthRoutes(r)
	routes.ViewRoutes(r)

	r.Run(":8080")
}
