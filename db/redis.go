package db

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:os.Getenv("REDIS_ADDR"),
		Username: "default",
		Password: os.Getenv("REDIS_KEY"),
		DB:       0,
	})
}
