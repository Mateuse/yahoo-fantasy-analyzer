package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

func ConnectToMySQL() error {
	sqlUser := os.Getenv("SQL_USER")
	sqlPassword := os.Getenv("SQL_PASSWORD")
	sqlHost := os.Getenv("SQL_HOST")
	sqlPort := os.Getenv("SQL_PORT")
	sqlDbName := os.Getenv("SQL_DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		sqlUser, sqlPassword, sqlHost, sqlPort, sqlDbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to mysql: %w", err)
	}

	log.Println("Connected to MySQL successfully")
	return nil
}

func AddRefreshToken(userId, refreshToken string) error {
	refreshTokenEntry := models.RefreshToken{UserId: userId, RefreshToken: refreshToken}
	if err := DB.Create(&refreshTokenEntry).Error; err != nil {
		return fmt.Errorf("failed to add refresh token for user: %w", err)
	}
	log.Printf("User %s Refresh Token added", userId)
	return nil
}

func GetRefreshToken(userId string) (string, error) {
	var refreshTokenEntry models.RefreshToken

	if err := DB.First(&refreshTokenEntry, "user_id = ?", userId).Error; err != nil {
		return "", fmt.Errorf("failed to get refresh token for user %s: $w", userId, err)
	}

	return refreshTokenEntry.RefreshToken, nil
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

type TeamWeeklyData struct {
	TeamID          string `gorm:"column:team_id"`
	Week            int    `gorm:"column:week"`
	ProjectedPoints string `gorm:"column:projected_points"`
	FinalPoints     string `gorm:"column:final_points"`
}

func AddTeamWeekData(teamId string, week int, projectedPoints, finalPoints string) error {
	data := TeamWeeklyData{
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

func SaveScheduleGameInDB(scheduleGame *models.ScheduleGame) error {
	err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // Match on the primary key `id`
		DoNothing: true,                          // Ignore the duplicate
	}).Create(scheduleGame).Error

	if err != nil {
		return fmt.Errorf("failed to insert or update data into DB: %w", err)
	}
	return nil
}

func GetTeamNextGameDB(teamAbbrev string) (*models.ScheduleGame, error) {
	var nextGame models.ScheduleGame

	err := DB.Where("(home_team_abbrev = ? OR away_team_abbrev = ?) AND start_time_utc > ?", teamAbbrev, teamAbbrev, time.Now()).
		Order("start_time_utc ASC").
		First(&nextGame).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no upcoming games found for team %s", teamAbbrev)
		}
		return nil, fmt.Errorf("failed to query next game: %s", err)
	}

	return &nextGame, nil
}

func SavePlayerStatsDB(player models.Player) error {
	statsJson, err := json.Marshal(player)

	if err != nil {
		return fmt.Errorf("failed to marshal player to JSON: %w", err)
	}

	nextGameTime, err := GetTeamNextGameDate(player.TeamAbbreviation)
	if err != nil {
		return fmt.Errorf("failed to get next game date: %w", err)
	}

	adjustedNextGameTime, err := AdjustTimePST(nextGameTime)
	if err != nil {
		return fmt.Errorf("failed to adjust next game time to PST: %w", err)
	}

	playerStats := models.PlayerStats{
		PlayerID:         player.PlayerKey,
		TeamAbbreviation: player.TeamAbbreviation,
		Stats:            string(statsJson),
		NextUpdate:       *adjustedNextGameTime,
	}

	err = DB.Save(&playerStats).Error
	if err != nil {
		return fmt.Errorf("failed to save player stats: %w", err)
	}

	return nil
}

func GetPlayerStatsDB(playerID string) (*models.PlayerStats, error) {
	var playerStats models.PlayerStats

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

func GetNHLPlayerByName

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
