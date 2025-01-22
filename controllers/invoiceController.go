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

var invoiceCollection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		results, err := invoiceCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can not reterive list of invoices"})
		}
		var allInvoices []bson.M
		if err = results.All(ctx, &allInvoices); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"invoices": allInvoices})
	}
}
func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		invoiceID := c.Param("invoice_id")
		var invoice models.Invoice
		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceID}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured"})
		}
		var invoiceView models.InvocieViewFormat
		allOrderItems, err := ItemsByOrder(invoice.OrderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		invoiceView.PaymentDueDate = invoice.PaymentDueDate
		invoiceView.PaymentMethod = nil
		if invoice.PaymentMethod != nil {
			invoiceView.PaymentMethod = invoice.PaymentMethod
		}
		invoiceView.InvoiceID = invoice.InvoiceID
		invoiceView.PaymentStatus = invoice.PaymentStatus
		invoiceView.PaymentDue = allOrderItems[0]["payment_due"]
		invoiceView.TableNumber = allOrderItems[0]["table_number"]
		invoiceView.OrderDetails = allOrderItems[0]["order_items"]
		c.JSON(http.StatusOK, invoiceView)
	}
}
func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		var order models.Order
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "details": "invalid invoice data in body data"})
			return
		}

		if invoice.PaymentStatus == nil {
			status := models.PaymentStatusPending
			invoice.PaymentStatus = &status
		}
		findErr := orderCollection.FindOne(ctx, bson.M{
			"order_id": invoice.OrderID,
		}).Decode(&order)
		if findErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order not found", "id": invoice.OrderID})
			return
		}
		invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.CreatedAt = helpers.RFC3339CurrentTime()
		invoice.UpdatedAt = helpers.RFC3339CurrentTime()
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceID = invoice.ID.Hex()
		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		result, insErr := foodCollection.InsertOne(ctx, invoice)
		if insErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		invoiceID := c.Param("invoice_id")
		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceID}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured"})
		}
		var updateObj primitive.D
		if invoice.PaymentMethod != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_method", Value: invoice.PaymentMethod})

		}
		if invoice.PaymentStatus == nil {
			status := models.PaymentStatusPending
			invoice.PaymentStatus = &status
		}
		updateObj = append(updateObj, bson.E{Key: "payment_status", Value: invoice.PaymentStatus})
		invoice.UpdatedAt = helpers.RFC3339CurrentTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})
		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		res, err := invoiceCollection.UpdateOne(ctx, bson.M{"invoice_id": invoiceID}, bson.D{
			{Key: "$set", Value: updateObj},
		}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)

	}
}
