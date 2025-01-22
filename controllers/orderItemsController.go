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
	"go.mongodb.org/mongo-driver/mongo"
)

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "oredrItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		res, err := orderItemCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var allOrderItems []bson.M
		if err = res.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"order_items": allOrderItems})
	}
}
func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param(("order_id"))
		allOrderItems, err := ItemsByOrder(orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"order_id": orderID, "items": allOrderItems})
	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	var items []primitive.M
	return items, nil
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderItemID := c.Param("order_item_id")
		var orderItem models.OrderItem
		err := orderItemCollection.FindOne(ctx, bson.M{"orderItem_id": orderItemID}).Decode(&orderItemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var OrderItemPack models.OrderItemPack
		var order models.Order
		if err := c.BindJSON(&OrderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		order.OrderDate = helpers.RFC3339CurrentTime()
		orderItemsToInsert := []interface{}{}
		order.TableID = OrderItemPack.TableID
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range OrderItemPack.OrderItems {
			orderItem.OrderID = order_id
			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": validationErr.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.UpdateAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.OrderItemId = orderItem.ID.Hex()
			var num = toFixed(*orderItem.UnitPrice, 2)
			orderItem.UnitPrice = &num
			orderItemsToInsert = append(orderItemsToInsert, orderItem)

		}
		res, err := orderItemCollection.InsertMany(ctx, orderItemsToInsert)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"inserted_items": res})
	}
}
