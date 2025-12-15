package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func Client(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDRESS"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	})
	return rdb
}
