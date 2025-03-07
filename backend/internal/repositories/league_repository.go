package repositories

import (
	"encoding/json"
	"fmt"
)

func SaveLeagueSettingsToDB(leagueId string, settings map[string]interface{}) error {
	jsonSettings, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to serialize settings to JSON: %w", err)
	}

	query := `
        INSERT INTO league_settings (league_id, settings)
        VALUES (?, ?)
        ON DUPLICATE KEY UPDATE settings = VALUES(settings), last_updated = NOW()
    `
	if err := DB.Exec(query, leagueId, jsonSettings).Error; err != nil {
		return fmt.Errorf("failed to save league settings to database: %w", err)
	}

	return nil
}

func GetLeagueSettingsFromDB(leagueId string) (map[string]interface{}, error) {
	var jsonSettings string
	query := `SELECT settings FROM league_settings WHERE league_id = ?`
	err := DB.Raw(query, leagueId).Scan(&jsonSettings).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch league settings from database: %w", err)
	}

	if jsonSettings == "" {
		return nil, nil // No settings found
	}

	var settings map[string]interface{}
	if err := json.Unmarshal([]byte(jsonSettings), &settings); err != nil {
		return nil, fmt.Errorf("failed to deserialize settings JSON: %w", err)
	}

	return settings, nil
}
