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
	
	routes.AuthRoutes(r)

	r.Run(":8080")
}
