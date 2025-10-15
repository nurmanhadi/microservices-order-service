package main

import (
	"order-service/config"
	"order-service/config/database"
	docs "order-service/docs"
	"order-service/pkg/env"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Order Service API
// @version 1.0
// @description This is the Order Service for managing orders and items.
// @BasePath /api
func main() {
	env.NewEnv()
	logger := config.NewLogger()
	db := database.NewSql()
	defer db.Close()

	validation := config.NewValidator()
	router := config.NewRouter()
	conn, ch := config.NewBroker()
	defer conn.Close()
	defer ch.Close()

	config.Setup(&config.DependenciesConfig{
		DB:         db,
		Logger:     logger,
		Validation: validation,
		Router:     router,
		Ch:         ch,
	})

	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := router.Run("0.0.0.0:8082")
	if err != nil {
		logger.WithError(err).Fatal("failed to start server")
	}
}
