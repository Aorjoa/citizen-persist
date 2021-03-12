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

// CitizenRequest should use to map incoming request
type CitizenRequest struct {
	CitizenID *string `json:"citizen_id"`
}

// ErrorMessageResponse should use to respose message as an application standard
type ErrorMessageResponse struct {
	Message string `json:"message"`
}

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
		var ci CitizenRequest
		if err := c.BodyParser(&ci); err != nil {
			logger.Error("unable to parse request")
			return c.Status(http.StatusBadRequest).JSON(&ErrorMessageResponse{Message: "unable to parse request"})
		}

		_, err := rdb.Get(ctx, *ci.CitizenID).Result()
		if err == nil {
			return c.SendStatus(http.StatusConflict)
		}

		err = rdb.Set(ctx, *ci.CitizenID, true, 10*time.Second).Err()
		if err != nil {
			logger.Error("unable to set citizen id to redis")
			return c.Status(http.StatusInternalServerError).JSON(&ErrorMessageResponse{Message: "unable to set citizen id to redis to parse request"})
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
