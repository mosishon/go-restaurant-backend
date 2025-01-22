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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		results, err := orderCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can not reterive list of orders"})
		}
		var allOrders []bson.M
		if err = results.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, bson.M{"orders": allOrders, "count": len(allOrders)})
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderID := c.Param("order_id")
		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured"})
		}
		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table
		var order models.Order
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "details": "invalid order data in body data"})
			return
		}

		validationErr := validate.Struct(order)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		if order.TableID != nil {
			findErr := tableCollection.FindOne(ctx, bson.M{
				"table_id": order.TableID,
			}).Decode(&table)
			if findErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "table not found", "id": order.TableID})
				return
			}

		}
		order.CreatedAt = helpers.RFC3339CurrentTime()
		order.UpdatedAt = helpers.RFC3339CurrentTime()
		order.ID = primitive.NewObjectID()
		order.OrderID = order.ID.Hex()
		result, insErr := orderCollection.InsertOne(ctx, order)
		if insErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created", "details": insErr.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table
		var order models.Order
		var updateObj primitive.D
		orderID := c.Param("order_id")
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if order.TableID != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "table was not found."})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.TableID})
		}
		order.UpdatedAt = helpers.RFC3339CurrentTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.UpdatedAt})
		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		res, err := orderCollection.UpdateOne(ctx, bson.M{"order_id": orderID}, bson.D{
			{Key: "$set", Value: updateObj},
		}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func OrderItemOrderCreator(order models.Order) string {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	order.CreatedAt = helpers.RFC3339CurrentTime()
	order.UpdatedAt = helpers.RFC3339CurrentTime()
	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()
	orderCollection.InsertOne(ctx, order)
	return order.OrderID
}
