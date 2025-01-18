package routes

import (
	"go-simple-shop/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	router.GET("/menus", controllers.GetMenus())
	router.GET("/menu/:menu_id", controllers.GetMenu())
	router.POST("/menus", controllers.CreateMenu())
	router.PATCH("/menus/:menu_id", controllers.UpdateMenu())
}
