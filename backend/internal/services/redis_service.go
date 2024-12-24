package services

import (
	"context"
	"os"
)

var ctx = context.Background()
var redisURL = os.Getenv("REDIS_URL")
var redisPassword = os.Getenv("REDIS")
