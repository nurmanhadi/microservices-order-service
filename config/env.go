package config

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

type envConfig struct {
	DB dbConfig
}

var ENV envConfig

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
	ENV.DB = envDB
}

func envAsInt(value string, defaultValue int) int {
	newValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return newValue
}
