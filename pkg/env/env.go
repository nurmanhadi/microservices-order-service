package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type dbConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	Name         string
	MaxIdleConns int
	MaxPoolConns int
	MaxLifetime  int
}
type apiConfig struct {
	Product string
}
type brokerConfig struct {
	VirtualHost string
	Host        string
	Port        string
	Username    string
	Password    string
	Exchange    struct {
		ProductService string
	}
	RouteKey struct {
		OrderCreated string
	}
}
type envConfig struct {
	DB     dbConfig
	API    apiConfig
	Broker brokerConfig
}

var CONF envConfig

func NewEnv() {
	godotenv.Load()

	envDB := dbConfig{
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		Username:     os.Getenv("DB_USERNAME"),
		Password:     os.Getenv("DB_PASSWORD"),
		Name:         os.Getenv("DB_DATABASE"),
		MaxIdleConns: envAsInt(os.Getenv("DB_MAX_IDLE_CONNS"), 5),
		MaxPoolConns: envAsInt(os.Getenv("DB_MAX_POOL_CONNS"), 10),
		MaxLifetime:  envAsInt(os.Getenv("DB_CONN_MAX_LIFETIME"), 300),
	}
	envAPI := apiConfig{
		Product: os.Getenv("API_PRODUCT"),
	}
	envBroker := brokerConfig{
		VirtualHost: os.Getenv("BROCKER_VIRTUAL_HOSTS"),
		Host:        os.Getenv("BROKER_HOST"),
		Port:        os.Getenv("BROKER_PORT"),
		Username:    os.Getenv("BROKER_USERNAME"),
		Password:    os.Getenv("BROKER_PASSWORD"),
		Exchange: struct{ ProductService string }{
			ProductService: os.Getenv("BROKER_EXCHANGE_PRODUCT_SERVICE"),
		},
		RouteKey: struct{ OrderCreated string }{
			OrderCreated: os.Getenv("BROCKER_ROUTE_ORDER_CREATED"),
		},
	}

	CONF.API = envAPI
	CONF.DB = envDB
	CONF.Broker = envBroker
}

func envAsInt(value string, defaultValue int) int {
	newValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return newValue
}
