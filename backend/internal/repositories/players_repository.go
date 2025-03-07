package repositories

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SavePlayerStatsDB(player models.Player) error {
	// 1. Get the next game time
	nextGameTime, err := GetTeamNextGameDate(player.TeamAbbreviation)
	if err != nil {
		return fmt.Errorf("failed to get next game date: %w", err)
	}

	// 2. Convert to PST
	adjustedNextGameTime, err := utils.AdjustTimePST(nextGameTime)
	if err != nil {
		return fmt.Errorf("failed to adjust next game time to PST: %w", err)
	}
	player.NextUpdate = *adjustedNextGameTime

	// 3. Convert EligiblePositions, Stats, and AdvancedStats to JSON strings
	statsJSON, err := json.Marshal(player.Stats)
	if err != nil {
		return fmt.Errorf("failed to marshal stats to JSON: %w", err)
	}

	advancedStatsJSON, err := json.Marshal(player.AdvancedStats)
	if err != nil {
		return fmt.Errorf("failed to marshal advanced stats to JSON: %w", err)
	}

	eligiblePositionsJSON, err := json.Marshal(player.EligiblePositions)
	if err != nil {
		return fmt.Errorf("failed to marshal eligible positions to JSON: %w", err)
	}

	// 4. Execute SQL query with ON DUPLICATE KEY UPDATE
	err = DB.Exec(`
        INSERT INTO players (
            player_key, player_id, full_name, first_name, last_name, ascii_first, ascii_last,
            team_full_name, team_abbr, team_url, uniform_number, display_position, 
            headshot_url, image_url, is_undroppable, position_type, 
            stats, advanced_stats, eligible_positions, next_update
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE 
            stats = VALUES(stats),
            advanced_stats = VALUES(advanced_stats),
            eligible_positions = VALUES(eligible_positions),
            next_update = VALUES(next_update);
    `, player.PlayerKey, player.PlayerID,
		player.Name.FullName, player.Name.FirstName, player.Name.LastName, player.Name.AsciiFirst, player.Name.AsciiLast,
		player.TeamFullName, player.TeamAbbreviation, player.TeamURL, player.UniformNumber, player.DisplayPosition,
		player.HeadshotURL, player.ImageURL, player.IsUndroppable, player.PositionType,
		string(statsJSON), string(advancedStatsJSON), string(eligiblePositionsJSON), player.NextUpdate).Error

	if err != nil {
		return fmt.Errorf("failed to save player stats: %w", err)
	}

	return nil
}

func GetPlayerStatsDB(playerID string) (*models.Player, error) {
	var playerStats models.Player

	err := DB.First(&playerStats, "player_id = ?", playerID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to fetch player stats: %w", err)
	}

	return &playerStats, nil
}

func SavePlayerIDMapping(yahooID string, nhlID string, name string, team string) error {
	playerMapping := models.PlayerIDMapping{
		YahooPlayerID: yahooID,
		NHLPlayerID:   nhlID,
		PlayerName:    name,
		TeamAbbr:      team,
	}

	err := DB.Save(&playerMapping).Error
	if err != nil {
		return fmt.Errorf("failed to save player ID mapping: %w", err)
	}
	return nil
}

func GetNHLPlayerID(yahooID string) (string, error) {
	var playerMapping models.PlayerIDMapping

	err := DB.First(&playerMapping, "yahoo_player_id = ?", yahooID).Error
	if err != nil {
		return "", fmt.Errorf("failed to find NHL Player ID for Yahoo ID %s: %w", yahooID, err)
	}

	return playerMapping.NHLPlayerID, nil
}

func SaveNhlPlayerToDB(players []*models.NHLPlayer) error {
	for _, player := range players {
		err := DB.Where("id = ?", player.ID).FirstOrCreate(player).Error
		if err != nil {
			return fmt.Errorf("failed to store player %d: %w", player.ID, err)
		}
	}

	return nil
}

func GetNhlPlayers() ([]models.NHLPlayer, error) {
	var nhlPlayers []models.NHLPlayer
	err := DB.Select("id, first_name, last_name").Find(&nhlPlayers).Error
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch NHL players: %v", err)
	}

	return nhlPlayers, nil
}

func GetYahooPlayers() ([]models.YahooPlayer, error) {
	var yahooPlayers []models.YahooPlayer
	err := DB.Select("id, full_name, team_name").Find(&yahooPlayers).Error
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch Yahoo players: %v", err)
	}
	return yahooPlayers, nil
}

func SaveYahooPlayerToDB(players []*models.YahooPlayer) error {
	for _, player := range players {
		err := DB.Where("id = ?", player.ID).FirstOrCreate(player).Error
		if err != nil {
			return fmt.Errorf("failed to store player %s: %w", player.ID, err)
		}
	}

	return nil
}

func SavePlayerIDMappingToDB(mappings []models.PlayerIDMapping) error {
	if err := DB.Create(&mappings).Error; err != nil {
		log.Fatalf("Failed to insert player mappings: %v", err)
		return err
	}

	return nil
}

func GetMappedPlayerByName(playerName string) (*models.PlayerIDMapping, error) {
	var player *models.PlayerIDMapping

	err := DB.Where("player_name = ?", playerName).First(&player).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch player %s: %w", playerName, err)
	}

	return player, nil
}

func GetYahooPlayerByName(playerName string) (*models.Player, error) {
	var player *models.Player

	err := DB.Where("name = ?", playerName).First(&player).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch player %s: %w", playerName, err)
	}

	return player, nil
}

func GetNhlPlayerById(playerId string) (*models.NHLPlayer, error) {
	var player *models.NHLPlayer

	err := DB.Where("id = ?", playerId).First(&player).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch player %s: %w", &playerId, err)
	}

	return player, nil
}

func SavePlayerGameStats(playerGameStats []*models.PlayerGameStat) error {
	err := DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(playerGameStats).Error

	if err != nil {
		return fmt.Errorf("failed to insert player game stats: %w", err)
	}

	log.Printf("Successfully inserted player game stats")
	return nil
}
