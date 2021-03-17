package main

import (
	"context"
	"os"

	"github.com/Aorjoa/citizen-persist/citizen"
	"github.com/Aorjoa/citizen-persist/model"
	"github.com/Aorjoa/citizen-persist/mq"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ctx = context.Background()

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "topic",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	r.SetOffset(kafka.LastOffset)

	m := mq.NewKafka(nil, r, ctx)

	cn := "host=localhost user=root password=root dbname=citizens port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(cn), &gorm.Config{})
	if err != nil {
		logger.Error("unable to connect database", zap.Error(err))
	}

	p := citizen.NewPersistent(logger, db)

	logger.Info("start consume message")
	for {
		k, v, err := m.ReadMessage()
		if err != nil {
			break
		}
		logger.Info("receive message", zap.ByteString("type", k), zap.ByteString("citizenID", v))
		id := string(v)
		mc := model.Citizen{
			CID: &id,
		}
		if err := p.Create(&mc); err != nil {
			logger.Error("unable to create citizen", zap.String("id", id), zap.Error(err))
		}
	}
}
