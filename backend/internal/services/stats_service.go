package services

import (
	"fmt"
	"strconv"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
)

func GetLeaguePlayerStats(statModifiers []models.StatModifier, player models.Player) (*models.Player, float64, error) {
	// Create a map of StatModifiers for quick lookup by StatID
	modifierMap := make(map[string]float64)
	for _, modifier := range statModifiers {
		modifierMap[modifier.StatID] = modifier.Value
	}

	totalPoints := 0.0

	// Adjust player stats
	for i, stat := range player.Stats {
		// Check if there is a modifier for the current StatID
		modifier, exists := modifierMap[stat.StatID]
		if exists {
			// Convert stat value to float for calculation
			statValue, err := strconv.ParseFloat(stat.Value, 64)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to parse stat value for StatID %s: %w", stat.StatID, err)
			}

			// Adjust the stat value
			adjustedValue := statValue * modifier

			// Update the player's stat with the adjusted value
			player.Stats[i].Value = fmt.Sprintf("%.2f", adjustedValue)

			// Add to the total points
			totalPoints += adjustedValue
		}
	}

	return &player, totalPoints, nil
}

func GetProjectedvsExpected(userSession, fTeamId string) ([]models.ProjectedVsActualStats, error) {

	return nil, nil
}
