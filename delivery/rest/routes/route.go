package routes

import (
	"order-service/delivery/rest/handler"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Router       *gin.Engine
	OrderHandler handler.OrderHandler
}

func (r *RouteConfig) Setup() {
	api := r.Router.Group("/api")

	order := api.Group("/orders")
	order.POST("/", r.OrderHandler.AddOrder)
	order.GET("/", r.OrderHandler.GetAllOrder)
	order.GET("/users/:id", r.OrderHandler.GetAllOrderByUserID)
	order.GET("/:id", r.OrderHandler.GetOrderByID)
}
