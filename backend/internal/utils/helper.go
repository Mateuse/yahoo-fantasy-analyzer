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

func GetStartOfCurrentWeek() time.Time {
	// Load EST location
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc)

	startOfWeek := now
	for startOfWeek.Weekday() != time.Monday {
		startOfWeek = startOfWeek.AddDate(0, 0, -1)
	}
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 9, 0, 0, 0, loc)

	return startOfWeek
}

func GenerateMatchupKey(teamA, teamB, week string) string {
	if teamA < teamB {
		return fmt.Sprintf("%s_vs_%s_%s", teamA, teamB, week)
	}
	return fmt.Sprintf("%s_vs_%s_%s", teamB, teamA, week)
}

func AdjustTimePST(givenTime *time.Time) (*time.Time, error) {
	if givenTime == nil {
		return nil, fmt.Errorf("given time is nil")
	}

	// Load PST timezone
	location, err := time.LoadLocation("America/Los_Angeles") // PST timezone
	if err != nil {
		return nil, fmt.Errorf("failed to load PST location: %w", err)
	}

	// Convert given time to PST
	timePST := givenTime.In(location)

	nextDay8AM := time.Date(
		timePST.Year(), timePST.Month(), timePST.Day()+1,
		8, 0, 0, 0,
		location,
	)

	return &nextDay8AM, nil
}
