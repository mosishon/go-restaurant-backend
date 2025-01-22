package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO write tags
type Order struct {
	ID        primitive.ObjectID
	OrderDate time.Time
	OrderID   string
	UpdatedAt time.Time
	CreatedAt time.Time
	TableID   *string
}
