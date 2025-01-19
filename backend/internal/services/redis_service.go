package services

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisURL = os.Getenv("REDIS_URL")
var redisPassword = os.Getenv("REDIS")

var redisClient = redis.NewClient(&redis.Options{
	Addr:     redisURL,
	Password: redisPassword,
	DB:       0,
})
