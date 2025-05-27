package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Amount    float64     `json:"amount"`
	Items     []OrderItem `json:"items"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}
type OrderItem struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
