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
		Addr:     "redis-18816.c264.ap-south-1-1.ec2.redns.redis-cloud.com:18816",
		Username: "default",
		Password: os.Getenv("REDIS_KEY"),
		DB:       0,
	})
}
