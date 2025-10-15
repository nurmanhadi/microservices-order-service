package database

import (
	"fmt"
	"log"
	"order-service/pkg/env"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewSql() *sqlx.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.CONF.DB.Host,
		env.CONF.DB.Port,
		env.CONF.DB.Username,
		env.CONF.DB.Password,
		env.CONF.DB.Name,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed connect to database: %s", err)
	}
	db.SetMaxIdleConns(env.CONF.DB.MaxIdleConns)
	db.SetMaxOpenConns(env.CONF.DB.MaxPoolConns)
	db.SetConnMaxLifetime(time.Duration(env.CONF.DB.MaxLifetime) * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed ping to database: %s", err)
	}
	return db
}
