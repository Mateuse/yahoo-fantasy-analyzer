package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"gorm.io/gorm"
)

func SaveTeamWeeklyStats(teamWeeklyStats []*models.TeamWeeklyStats) error {
	// Create a slice to hold the JSON-encoded stats
	type TeamWeeklyStatsWithJSON struct {
		models.TeamWeeklyStats
		StatsJSON string
	}

	var teamWeeklyStatsWithJSON []TeamWeeklyStatsWithJSON

	// Convert Stats to JSON string for each entry
	for _, stats := range teamWeeklyStats {
		statsJSON, err := json.Marshal(stats.Stats)
		if err != nil {
			return fmt.Errorf("failed to marshal stats to JSON: %w", err)
		}

		teamWeeklyStatsWithJSON = append(teamWeeklyStatsWithJSON, TeamWeeklyStatsWithJSON{
			TeamWeeklyStats: *stats,
			StatsJSON:       string(statsJSON),
		})
	}

	// Perform bulk insert with ON DUPLICATE KEY UPDATE
	err := DB.Transaction(func(tx *gorm.DB) error {
		for _, stats := range teamWeeklyStatsWithJSON {
			err := tx.Exec(`
                INSERT INTO team_weekly_stats (
                    id, week, team_key, stats, points
                ) VALUES (?, ?, ?, ?, ?)
                ON DUPLICATE KEY UPDATE 
                    stats = VALUES(stats),
                    points = VALUES(points);
            `, stats.ID, stats.Week, stats.TeamKey, stats.StatsJSON, stats.Points).Error

			if err != nil {
				return fmt.Errorf("failed to insert team weekly stats: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to save team weekly stats: %w", err)
	}

	return nil
}
