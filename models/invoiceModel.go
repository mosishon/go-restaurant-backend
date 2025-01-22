package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentMethodType string
type PaymentStatusType string

const (
	PaymentMethodCard PaymentMethodType = "CARD"
	PaymentMethodCash PaymentMethodType = "CASH"

	PaymentStatusPending PaymentStatusType = "PENDING"
	PaymentStatusPaid    PaymentStatusType = "PAID"
)

type Invoice struct {
	ID             primitive.ObjectID `bson:"_id"`
	InvoiceID      string             `json:"invoice_id"`
	OrderID        string             `json:"order_id"`
	PaymentMethod  *PaymentMethodType `json:"payment_method"`
	PaymentStatus  *PaymentStatusType `json:"payment_status" validate:"required"`
	PaymentDueDate time.Time          `json:"payment_due_date"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

// TODO write tags

type InvocieViewFormat struct {
	InvoiceID      string
	PaymentMethod  *PaymentMethodType
	orderID        string
	PaymentStatus  *PaymentStatusType
	PaymentDue     interface{}
	TableNumber    interface{}
	PaymentDueDate time.Time
	OrderDetails   interface{}
}
