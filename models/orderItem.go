package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID          primitive.ObjectID
	Quantity    *string
	UnitPrice   *float64
	CreatedAt   time.Time
	UpdateAt    time.Time
	FoodID      *string
	OrderItemId string
	OrderID     string
}

type OrderItemPack struct {
	TableID    *string
	OrderItems []OrderItem
}
