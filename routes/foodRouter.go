package routes

import (
	"go-simple-shop/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRouter(router *gin.Engine) {
	router.GET("/foods", controllers.GetFoods())
	router.GET("/food/:food_id", controllers.GetFood())
	router.POST("/foods", controllers.CreateFood())
	router.PATCH("/foods/:food_id", controllers.UpdateFood())

}
