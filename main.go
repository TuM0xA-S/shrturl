package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
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

	shrt := NewShrt(ctx, ShrtCfg{BaseURL: "http://localhost:1323/", DB: db})

	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	shrt.Setup(e)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
