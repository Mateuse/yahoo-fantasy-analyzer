package utils

import (
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

func GetInt(data map[string]interface{}, key string) int {
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
