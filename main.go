package main

import (
	"context"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4/middleware"
)

type config struct {
	RedisAddr     string        `env:"REDIS_ADDR"`
	RedisUser     string        `env:"REDIS_USER"`
	RedisPassword string        `env:"REDIS_PASSWORD"`
	RedisDB       int           `env:"REDIS_DB"`
	ExpiryTime    time.Duration `env:"EXPIRY_TIME"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	db := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Username: cfg.RedisUser,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	shrt := NewShrt(ctx, ShrtCfg{DB: db, ExpirationTime: cfg.ExpiryTime})
	shrt.Use(middleware.Logger(), middleware.Recover())

	// Start server
	shrt.Logger.Fatal(shrt.Start(":1323"))
}
