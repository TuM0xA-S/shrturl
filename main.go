package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()

	db := redis.NewClient(&redis.Options{
		Addr:     "",
		Username: "",
		Password: "",
		DB:       0,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	shrt := NewShrt(ctx, ShrtCfg{DB: db, ExpirationTime: 10 * time.Second})
	shrt.Use(middleware.Logger(), middleware.Recover())

	// Start server
	shrt.Logger.Fatal(shrt.Start(":1323"))
}
