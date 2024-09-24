package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	OrderPending    OrderStatus = "OrderPending"
	OrderProcessing OrderStatus = "OrderProcessing"
	OrderCompleted  OrderStatus = "OrderCompleted"
	OrderCancelled  OrderStatus = "OrderCancelled"
	OrderDelivered  OrderStatus = "OrderDelivered"
)

type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"orderId"`
	Version     int64              `json:"version" bson:"version"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	Products    []Product          `json:"products" bson:"products"`
	User        string             `json:"user" bson:"user"`
	TotalAmount float64            `json:"totalAmount" bson:"totalAmount"`
	Status      OrderStatus        `json:"status" bson:"status"`
	Updates     []OrderUpdate      `json:"updates" bson:"updates"`
}

type Product struct {
	Name     string    `json:"name" bson:"name"`
	UpdateAt time.Time `json:"updateAt" bson:"updateAt"`
	Price    float64   `json:"price" bson:"price"`
	Status   string    `json:"status" bson:"status"`
	Remarks  string    `json:"remarks" bson:"remarks"`
	Quantity uint64    `json:"quantity"`
}

type OrderUpdate struct {
	UpdateAt time.Time `json:"updateAt" bson:"updateAt"`
	Notes    string    `json:"notes" bson:"notes"`
	HandleBy string    `json:"handleBy" bson:"handleBy"`
}
