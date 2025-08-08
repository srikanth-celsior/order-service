package main

import (
	"log"
	"orders-service/database"
	"orders-service/pubsub"
	"orders-service/routers"
	"orders-service/utils"
	"os"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func main() {
	_ = godotenv.Load()
	utils.InitRedis()
	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	if err := pubsub.InitPubSub(); err != nil {
		log.Fatal("PubSub init failed:", err)
	}
	app := iris.New()
	routers.RegisterRoutes(app)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	_ = app.Listen(":" + port)
}
