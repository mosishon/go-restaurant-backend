package controllers

import (
	"context"
	"go-simple-shop/database"
	"go-simple-shop/helpers"
	"go-simple-shop/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection = database.OpenCollection(database.Client, "menu")
var validate = validator.New()

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		results, err := menuCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can not reterive list of menus"})
		}
		var allMenus []bson.M
		if err = results.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"menus": allMenus, "count": len(allMenus)})

	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		menuID := c.Param("menu_id")
		var menu models.Menu
		err := foodCollection.FindOne(ctx, bson.M{"menu_id": menuID}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured" + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"menu": menu})
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.MenuID = menu.ID.Hex()
		_, insErr := menuCollection.InsertOne(ctx, menu)
		if insErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		menuID := c.Param("menu_id")
		filter := bson.M{"menu_id": menuID}

		var updateObj primitive.D
		if menu.StartDate != nil && menu.EndDate != nil {
			if !helpers.InTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "kindly retype the time"})
				return
			}
		}
		updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.StartDate})
		updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.EndDate})

		if menu.Name != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
		}
		if menu.Category != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Category})
		}

		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.UpdatedAt})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		res, err := menuCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)

	}
}
