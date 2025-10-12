package database

import (
	"fmt"
	"log"
	"order-service/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewSql() *sqlx.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.ENV.DB.Host,
		config.ENV.DB.Port,
		config.ENV.DB.Username,
		config.ENV.DB.Password,
		config.ENV.DB.Name,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed connect to database: %s", err)
	}
	db.SetMaxIdleConns(config.ENV.DB.MaxIdleConns)
	db.SetMaxOpenConns(config.ENV.DB.MaxPoolConns)
	db.SetConnMaxLifetime(time.Duration(config.ENV.DB.MaxLifetime) * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed ping to database: %s", err)
	}
	return db
}
