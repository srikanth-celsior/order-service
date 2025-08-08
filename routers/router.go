package routers

import (
	"orders-service/handlers"
	"orders-service/middleware"

	"github.com/kataras/iris/v12"
)

func RegisterRoutes(app *iris.Application) {
	orders := app.Party("/orders", middleware.JWTMiddleware)
	{
		orders.Post("/", handlers.CreateOrder)
		orders.Get("/{id}", handlers.GetOrder)
		orders.Patch("/{id}/status", handlers.UpdateOrderStatus)
	}
}
