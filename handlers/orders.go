package handlers

import (
	"context"
	"encoding/json"
	"log"
	"orders-service/database"
	"orders-service/models"
	"orders-service/pubsub"
	"orders-service/utils"
	"time"

	"github.com/gofrs/uuid"
	"github.com/kataras/iris/v12"
)

func CreateOrder(ctx iris.Context) {
	var req models.Order

	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StopWithStatus(400)
		return
	}

	// userID := ctx.Values().GetString("user_id")

	row := database.DB.QueryRow(`
		INSERT INTO orders (user_id, amount, status, created_at)
		VALUES ($1, $2, 'pending', $3) RETURNING id`,
		req.UserID, req.Amount, time.Now())

	var orderID uuid.UUID
	if err := row.Scan(&orderID); err != nil {
		log.Println("Error inserting order:", err)
		ctx.StopWithStatus(500)
		return
	}

	for _, item := range req.Items {
		row := database.DB.QueryRow(`
			INSERT INTO order_items (order_id, product_id, product_name, quantity, price)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			orderID, item.ProductID, item.ProductName, item.Quantity, item.Price)
		var itemID uuid.UUID
		if err := row.Scan(&itemID); err != nil {
			log.Println("Error inserting order item:", err)
			ctx.StopWithStatus(500)
			return
		}
	}

	// Publish to PubSub
	_ = pubsub.PublishOrderEvent(map[string]interface{}{
		"order_id": orderID,
		"user_id":  req.UserID,
		"amount":   req.Amount,
	})

	ctx.JSON(iris.Map{
		"order_id": orderID,
		"status":   "pending",
	})
}

func GetOrder(ctx iris.Context) {
	id := ctx.Params().Get("id")
	rdb := utils.GetRedisClient()
	ctxRedis := context.Background()
	cacheKey := "order:" + id

	var o models.Order
	// Try to get from Redis first
	orderJson, err := rdb.Get(ctxRedis, cacheKey).Result()
	if err == nil {
		// Cache hit
		ctx.ContentType("application/json")
		ctx.Write([]byte(orderJson))
		return
	}

	// Cache miss, fetch from DB
	row := database.DB.QueryRow("SELECT id, user_id, amount, status, created_at FROM orders WHERE id = $1", id)
	err = row.Scan(&o.ID, &o.UserID, &o.Amount, &o.Status, &o.CreatedAt)
	if err != nil {
		ctx.StopWithStatus(404)
		return
	}

	// Store in Redis
	orderBytes, _ := json.Marshal(o)
	rdb.Set(ctxRedis, cacheKey, orderBytes, 300*time.Second)

	ctx.JSON(o)
}

func UpdateOrderStatus(ctx iris.Context) {
	id := ctx.Params().Get("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StopWithStatus(400)
		return
	}

	_, err := database.DB.Exec("UPDATE orders SET status = $1 WHERE id = $2", req.Status, id)
	if err != nil {
		ctx.StopWithStatus(500)
		return
	}

	ctx.JSON(iris.Map{"updated": true})
}
