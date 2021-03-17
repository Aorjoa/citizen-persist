package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/Aorjoa/citizen-persist/citizen"
	"github.com/Aorjoa/citizen-persist/middleware"
	"github.com/Aorjoa/citizen-persist/mq"
	redisStore "github.com/Aorjoa/citizen-persist/redis"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var ctx = context.Background()

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "topic",
		Balancer: &kafka.LeastBytes{},
	}

	m := mq.NewKafka(w, nil, ctx)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	rs := redisStore.NewStorage(rdb, ctx)

	c := citizen.NewHandler(logger, m, rs)

	app := fiber.New()

	ra := app.Group("/api")
	ra.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	v1 := app.Group("/api/v1")
	mwl := middleware.Logger(logger)
	v1.Use(mwl.LogWithContext)
	v1.Post("/citizens", c.PutCitizenIDToQueue)

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
