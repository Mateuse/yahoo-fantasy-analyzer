package repositories

import (
	"fmt"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"gorm.io/gorm"
)

func AddTeamWeekData(teamId string, week int, projectedPoints, finalPoints string) error {
	data := models.TeamWeeklyData{
		TeamID:          teamId,
		Week:            week,
		ProjectedPoints: projectedPoints,
		FinalPoints:     finalPoints,
	}

	err := DB.Create(&data).Error
	if err != nil {
		return fmt.Errorf("failed to insert data for team %s week %d: %w", teamId, week, err)
	}

	return nil
}

func GetTeamWeekData(teamId string, week int) (string, string, error) {
	var result struct {
		ProjectedPoints string
		FinalPoints     string
	}

	err := DB.Table("team_weekly_data").
		Select("projected_points, final_points").
		Where("team_id = ? AND week = ?", teamId, week).
		Scan(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", "", nil // No data found
		}
		return "", "", fmt.Errorf("failed to get data for team %s week %d: %w", teamId, week, err)
	}

	return result.ProjectedPoints, result.FinalPoints, nil
}
