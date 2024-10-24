package external

import "github.com/derickit/go-rest-api/internal/models/data"

type APIError struct {
	HTTPStatusCode int    `json:"httpStatusCode"`
	Message        string `json:"message"`
	DebugID        string `json:"debugId"`
	ErrorCode      string `json:"errorCode"`
}

type OrderInPut struct {
	Products []ProductInput `json:"products" binding:"required"`
}

type ProductInput struct {
	Name     string  `json:"name" bingding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Quantity uint64  `json:"quantity" binding:"required"`
}

type Order struct {
	ID          string             `json:"orderId"`
	Version     int64              `json:"version"`
	CreatedAt   string             `json:"createdAt"`
	UpdatedAt   string             `json:"updatedAt"`
	Products    []data.Product     `json:"products"`
	User        string             `json:"user"`
	TotalAmount float64            `json:"totalAmount"`
	Status      data.OrderStatus   `json:"status"`
	Updates     []data.OrderUpdate `json:"updates"`
}
