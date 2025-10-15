package config

import (
	"order-service/delivery/messaging/producer"
	"order-service/delivery/rest/handler"
	"order-service/delivery/rest/routes"
	"order-service/internal/repository"
	"order-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type DependenciesConfig struct {
	DB         *sqlx.DB
	Logger     *logrus.Logger
	Validation *validator.Validate
	Router     *gin.Engine
	Ch         *amqp091.Channel
}

func Setup(deps *DependenciesConfig) {
	// producer
	orderProd := producer.NewOrderProducer(deps.Ch)

	// repository
	orderRepo := repository.NewOrderRepository(deps.DB)

	// service
	orderServ := service.NewOrderService(deps.Logger, deps.Validation, orderRepo, orderProd)

	// handler
	orderHand := handler.NewOrderHandler(orderServ)

	// route
	route := &routes.RouteConfig{
		Router:       deps.Router,
		OrderHandler: orderHand,
	}
	route.Setup()
}
