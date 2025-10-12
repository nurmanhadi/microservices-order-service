package config

import (
	"order-service/delivery/rest/handler"
	"order-service/delivery/rest/routes"
	"order-service/internal/repository"
	"order-service/internal/service"
	"order-service/pkg/api"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type DependenciesConfig struct {
	DB         *sqlx.DB
	Logger     *logrus.Logger
	Validation *validator.Validate
	Router     *gin.Engine
}

func Setup(deps *DependenciesConfig) {
	// api client
	productAPI := api.NewProductAPI(ENV.API.Product)

	// repository
	orderRepo := repository.NewOrderRepository(deps.DB)

	// service
	orderServ := service.NewOrderService(deps.Logger, deps.Validation, orderRepo, productAPI)

	// handler
	orderHand := handler.NewOrderHandler(orderServ)

	// route
	route := &routes.RouteConfig{
		Router:       deps.Router,
		OrderHandler: orderHand,
	}
	route.Setup()
}
