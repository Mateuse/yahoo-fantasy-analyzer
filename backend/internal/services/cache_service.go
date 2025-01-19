package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func CacheResponse(sessionId string, requestType string, responseBody map[string]interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s", sessionId, requestType)

	data, err := json.Marshal(responseBody)
	if err != nil {
		return fmt.Errorf("failed to serialize body: %w", err)
	}

	err = redisClient.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to cache the response: %w", err)
	}

	return nil
}

func GetCachedResponse(sessionId string, requestType string) (map[string]interface{}, error) {
	key := fmt.Sprintf("%s:%s", sessionId, requestType)

	// Fetch data from Redis
	result, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// Not in cache
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch cached response: %w", err)
	}

	// Deserialize the JSON string into a map
	var responseBody map[string]interface{}
	if err := json.Unmarshal([]byte(result), &responseBody); err != nil {
		return nil, fmt.Errorf("failed to deserialize cached response: %w", err)
	}

	return responseBody, nil
}
