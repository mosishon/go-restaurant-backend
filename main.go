package main

import (
	"go-simple-shop/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	routes.FoodRouter(router)
	routes.OrderRoutes(router)
	routes.MenuRoutes(router)
	router.Run(":" + port)
}
