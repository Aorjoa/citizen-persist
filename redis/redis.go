package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type DBStorage interface {
	GetData(key string) (interface{}, error)
	SetData(key string, value interface{}, expiration time.Duration) error
}

type Storage struct {
	Client  *redis.Client
	Context context.Context
}

func NewStorage(redis *redis.Client, context context.Context) DBStorage {
	return &Storage{
		Client:  redis,
		Context: context,
	}
}

func (rs Storage) GetData(key string) (interface{}, error) {
	return rs.Client.Get(rs.Context, key).Result()
}

func (rs Storage) SetData(key string, value interface{}, expiration time.Duration) error {
	return rs.Client.Set(rs.Context, key, value, expiration).Err()
}
