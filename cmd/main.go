package main

import (
	"backend/config"
	"backend/routes"
	"backend/seeder"
	"backend/models"

	"encoding/json"
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()
	models.Migrate()

	err := seeder.AdminSeeder(config.DB)
	if err != nil {
		log.Fatalf("admin seeder is failed: %v", err)
	}

	r := gin.Default()

	// ✅ ADD THIS BLOCK
	r.SetFuncMap(template.FuncMap{
		"json": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
	})

	// ✅ MUST come AFTER SetFuncMap
	r.LoadHTMLGlob("templates/*")

	// Routes
	routes.AuthRoutes(r)
	routes.ProductRoutes(r)
	routes.CartRoutes(r)
	routes.WishlistRouts(r)
	routes.OrderRoutes(r)
	routes.UserProfileRoutes(r)
	routes.AdminRoutes(r)
	routes.ViewRoutes(r)

	r.Run(":8080")
}
