package controllers

import (
	"context"
	"go-simple-shop/database"
	"go-simple-shop/models"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection = database.OpenCollection(database.Client, "food")

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		results, err := foodCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can not reterive list of foods"})
		}
		var allFoods []bson.M
		if err = results.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, bson.M{"foods": allFoods, "count": len(allFoods)})
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		foodID := c.Param("food_id")
		var food models.Food
		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodID}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured"})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var food models.Food
		var menu models.Menu

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "details": "invalid food data in body data"})
			return
		}

		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		findErr := menuCollection.FindOne(ctx, bson.M{
			"menu_id": food.MenuID,
		}).Decode(&menu)
		if findErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "menu not found", "id": food.MenuID})
			return
		}
		food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.FoodID = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num
		result, insErr := foodCollection.InsertOne(ctx, food)
		if insErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu
		var food models.Food
		foodID := c.Param("food_id")
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var updateObj primitive.D
		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
		}
		if food.Price != nil {
			updateObj = append(updateObj, bson.E{Key: "price", Value: food.Price})
		}
		if food.FoodImage != nil {
			updateObj = append(updateObj, bson.E{Key: "food_image", Value: food.FoodImage})
		}
		if food.MenuID != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuID}).Decode(&menu)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "menu was not found."})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "menu_id", Value: food.MenuID})
		}
		food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.UpdatedAt})
		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		res, err := menuCollection.UpdateOne(ctx, bson.M{"food_id": foodID}, bson.D{
			{Key: "$set", Value: updateObj},
		}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	out := math.Pow(10, float64(precision))
	return float64(round(num*out)) / out
}
