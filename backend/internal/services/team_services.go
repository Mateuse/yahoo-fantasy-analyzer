package services

import (
	"fmt"
	"strconv"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func GetTeamWeeklyStats(sessionId, teamId string) (map[int]map[string]string, error) {

	leagueId, err := utils.TeamtoLeagueId(teamId)
	if err != nil {
		return nil, err
	}

	currentWeekStr, err := GetLeagueSetting(sessionId, leagueId, "CurrentWeek")

	if err != nil {
		return nil, err
	}

	if currentWeekStr == "1" {
		return nil, nil
	}

	currentWeek, err := strconv.Atoi(currentWeekStr)
	if err != nil {
		return nil, err
	}

	// Initialize a map to store the projected points for each week
	weeklyStats := make(map[int]map[string]string)

	// Loop through each week from 1 to currentWeek
	for week := 1; week <= currentWeek; week++ {
		// Check if data exists in the database
		projectedPoints, finalPoints, err := GetTeamWeekData(teamId, week)
		if err != nil {
			return nil, err
		}

		if projectedPoints != "" && finalPoints != "" {
			// Use existing data
			weeklyStats[week] = map[string]string{
				"expectedPoints": projectedPoints,
				"finalPoints":    finalPoints,
			}
			continue
		}

		// Fetch new data from Yahoo API if not found
		teamWeekly, err := GetTeamWeekStats(sessionId, teamId, strconv.Itoa(week))
		if err != nil {
			return nil, fmt.Errorf("failed to get team stats for week %d: %w", week, err)
		}

		// Store the fetched data in the map
		weeklyStats[week] = map[string]string{
			"expectedPoints": teamWeekly.ProjectedPoints,
			"finalPoints":    teamWeekly.FinalPoints,
		}

		// Insert new data into the database
		err = AddTeamWeekData(teamId, week, teamWeekly.ProjectedPoints, teamWeekly.FinalPoints)
		if err != nil {
			return nil, err
		}
	}

	return weeklyStats, nil
}
