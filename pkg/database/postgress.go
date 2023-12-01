package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"cekkustomer.com/configs"

	"github.com/sethvargo/go-envconfig"
)

var DB *sql.DB

func Init() (*sql.DB, error) {
	ctx := context.Background()
	var config configs.AppConfiguration

	if err := envconfig.Process(ctx, &config); err != nil {
		log.Fatal(err.Error())
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.EndPoint,
	)
	db, err := sql.Open(config.Database.Type, connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(70)
	db.SetMaxIdleConns(70)
	duration, err := time.ParseDuration("15m")
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
