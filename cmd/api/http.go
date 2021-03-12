package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var ctx = context.Background()

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	app := fiber.New()

	ra := app.Group("/api")
	ra.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	v1 := app.Group("/api/v1")
	v1.Post("/citizens", func(c *fiber.Ctx) error {
		_, err := rdb.Get(ctx, "cID").Result()
		if err == nil {
			c.SendStatus(http.StatusConflict)
		}

		err = rdb.Set(ctx, "cID", "1411900101002", 10*time.Second).Err()
		if err != nil {
			panic(err)
		}

		return c.SendString("Hello, World 👋!")
	})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		_ = <-sc
		logger.Info("gracefully shutting down")
		_ = app.Shutdown()
	}()

	logger.Info("server start listening")
	err = app.Listen(":3000")
	if err != nil {
		logger.Panic("server stop listening")
	}
	logger.Info("cleanup task")
}
