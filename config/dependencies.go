package config

import (
	"fmt"
	"log"
	"order-service/delivery/rest/middleware"
	"order-service/pkg/env"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}
func NewValidator() *validator.Validate {
	validator := validator.New()
	return validator
}
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ErrorHandling())
	return r
}
func NewBroker() (*amqp.Connection, *amqp.Channel) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", env.CONF.Broker.Username, env.CONF.Broker.Password, env.CONF.Broker.Host, env.CONF.Broker.Port, env.CONF.Broker.VirtualHost)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("failed to connect broker: %s", err.Error())
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel broker: %s", err.Error())
	}
	return conn, ch
}
