package main

import (
	"log"
	"orders-service/database"
	"orders-service/handlers"
	"orders-service/middleware"
	"orders-service/models"
	"orders-service/pubsub"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func main() {
	_ = godotenv.Load()
	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	if err := pubsub.InitPubSub(); err != nil {
		log.Fatal("PubSub init failed:", err)
	}

	app := iris.New()
	app.Post("/login", func(ctx iris.Context) {
		var req struct {
			UserID string `json:"user_id"`
		}
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StopWithStatus(400)
			return
		}

		token, err := GenerateJWT(req.UserID)
		if err != nil {
			ctx.StopWithStatus(500)
			return
		}

		ctx.JSON(iris.Map{
			"token": token,
		})
	})
	orders := app.Party("/orders", middleware.JWTMiddleware)
	{
		orders.Post("/", handlers.CreateOrder)
		orders.Get("/{id}", handlers.GetOrder)
		orders.Patch("/{id}/status", handlers.UpdateOrderStatus)
	}

	_ = app.Listen(":3000")
}
func GenerateJWT(userId string) (string, error) {
	claims := models.CustomClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 96)), // Expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecretKey := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(jwtSecretKey))
}
