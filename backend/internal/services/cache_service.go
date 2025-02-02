package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func CacheResponse(dataid string, requestType string, responseBody interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s", dataid, requestType)

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

func GetCachedResponse(dataid string, requestType string) (interface{}, error) {
	key := fmt.Sprintf("%s:%s", dataid, requestType)

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

func ClearCache(operation string) error {
	if operation == "" {
		//Flush all keys in redis
		err := redisClient.FlushAll(ctx).Err()
		if err != nil {
			return fmt.Errorf("failed to flush Redis cache: %w", err)
		}
		return nil
	}

	iter := redisClient.Scan(ctx, 0, fmt.Sprintf("*:%s", operation), 0).Iterator()
	for iter.Next(ctx) {
		err := redisClient.Del(ctx, iter.Val()).Err()
		if err != nil {
			return fmt.Errorf("failed to delet cache key %s: %w", iter.Val(), err)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error while scanning Redis keys: %w", err)
	}

	return nil
}
