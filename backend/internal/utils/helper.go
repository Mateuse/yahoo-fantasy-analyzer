package utils

import (
	"fmt"
	"strconv"
	"time"
)

// Helper functions for parsing values
func GetString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func GetInt64(data map[string]interface{}, key string) int64 {
	if val, ok := data[key].(string); ok {
		if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
			return intVal
		}
	}
	return 0
}

func GetInt(data map[string]interface{}, key string) int {
	// Handle numbers stored as float64
	if val, ok := data[key].(float64); ok {
		return int(val)
	}

	// Handle numbers stored as strings
	if val, ok := data[key].(string); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}

	return 0
}

func GetFloat(data map[string]interface{}, key string) float64 {
	if val, ok := data[key].(string); ok {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return 0.0
}

func ParseDate(value string) time.Time {
	if t, err := time.Parse("2006-01-02", value); err == nil {
		return t
	}
	return time.Time{}
}

func GetBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(string); ok {
		return val == "1"
	}
	return false
}

func GetTTL() time.Duration {
	now := time.Now()
	expiry := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return time.Until(expiry)
}

func GetCurrentNhlSeason() string {
	now := time.Now()
	year := now.Year()

	// NHL season starts in October and spans two years
	if now.Month() >= time.July {
		return fmt.Sprintf("%d%d", year, year+1)
	}
	return fmt.Sprintf("%d%d", year-1, year)
}
