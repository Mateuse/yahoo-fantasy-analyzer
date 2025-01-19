package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

type Session struct {
	UserId      string `json:"user_id`
	AccessToken string `json:"access_token`
	ExpiryTime  string `json:"expiry_time"`
}

func CreateUserSession(userId string, accessToken string, expiresIn float64) error {
	expiryTime := time.Now().Add(time.Second * time.Duration(expiresIn))

	sessionData := Session{
		UserId:      userId,
		AccessToken: accessToken,
		ExpiryTime:  expiryTime.Format(time.RFC3339),
	}

	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to encode session data: %w", err)
	}

	// Save the session in Redis
	err = redisClient.Set(ctx, userId, sessionJSON, time.Second*time.Duration(expiresIn)).Err()
	if err != nil {
		return fmt.Errorf("failed to save session to Redis: %w", err)
	}

	return nil
}

func GetUserSession(userId string) (*Session, error) {
	// Fetch session data from Redis
	sessionJSON, err := redisClient.Get(ctx, userId).Result()
	if err == redis.Nil {
		return nil, utils.NewNotFoundError(fmt.Sprintf("no session found for user: %s", userId))
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch session from Redis: %w", err)
	}

	// Decode JSON into Session struct
	var session Session
	err = json.Unmarshal([]byte(sessionJSON), &session)
	if err != nil {
		return nil, fmt.Errorf("failed to decode session data: %w", err)
	}

	return &session, nil
}
