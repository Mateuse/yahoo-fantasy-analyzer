package services

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/repositories"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
	"gorm.io/gorm"
)

func GetUserLeagues(userSession string) (map[string]interface{}, error) {

	url := "https://fantasysports.yahooapis.com/fantasy/v2/users;use_login=1/games/leagues"

	leaguesResponse, err := AuthHttpXMLRequest(userSession, url)
	if err != nil {
		return nil, err
	}

	leagues, err := ExtractLeaguesFromResponse(leaguesResponse)
	if err != nil {
		return nil, err
	}

	return leagues, nil
}

func GetLeague(userSession, leagueId string) (map[string]interface{}, error) {

	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/league/%s", leagueId)

	leagueResponse, err := AuthHttpXMLRequest(userSession, url)
	if err != nil {
		return nil, err
	}
	return leagueResponse, nil
}

func GetLeagueSettings(userSession, leagueId string) (*models.League, error) {
	// Check the database for existing league settings
	dbSettings, err := repositories.GetLeagueSettingsFromDB(leagueId)
	if err != nil {
		return nil, fmt.Errorf("error fetching league settings from database: %w", err)
	}

	if dbSettings != nil {
		// Convert DB settings to League object
		cachedSettings, err := MapToLeague(dbSettings)
		if err != nil {
			return nil, fmt.Errorf("error converting DB settings to league object: %w", err)
		}

		startOfCurrentWeek := utils.GetStartOfCurrentWeek()

		//If LastUpdated is after the start of the current week, return cached setting

		if cachedSettings.LastUpdated.After(startOfCurrentWeek) {
			return cachedSettings, nil
		}
	}

	// Make API call if not in cache
	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/league/%s/settings", leagueId)
	leagueSettingsResponse, err := AuthHttpXMLRequest(userSession, url)
	if err != nil {
		return nil, fmt.Errorf("error fetching league settings from API: %w", err)
	}

	// Convert API response to League object
	leagueSettings, err := MapToLeague(leagueSettingsResponse)
	if err != nil {
		return nil, fmt.Errorf("error converting API response to league: %w", err)
	}

	// Save settings to the database for future use
	err = repositories.SaveLeagueSettingsToDB(leagueId, leagueSettingsResponse)
	if err != nil {
		return nil, fmt.Errorf("error saving league settings to database: %w", err)
	}

	return leagueSettings, nil
}

func GetLeagueSetting(userSession, leagueId, setting string) (string, error) {
	leagueSettings, err := GetLeagueSettings(userSession, leagueId)
	if err != nil {
		return "", err
	}

	// Get the reflection value of the League struct
	leagueValue := reflect.ValueOf(leagueSettings)

	// Dereference the pointer if necessary
	if leagueValue.Kind() == reflect.Ptr {
		leagueValue = leagueValue.Elem()
	}

	// Check if it's a struct after dereferencing
	if leagueValue.Kind() != reflect.Struct {
		return "", errors.New("leagueSettings is not a struct")
	}

	// Get the field by name (case-sensitive)
	field := leagueValue.FieldByName(setting)
	if !field.IsValid() {
		return "", fmt.Errorf("setting %s not found in league settings", setting)
	}

	// Convert the field value to a string based on its kind
	switch field.Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Int, reflect.Int64:
		return strconv.Itoa(int(field.Int())), nil
	default:
		return "", errors.New("unsupported field type")
	}
}

func GetTeamWeekStats(sessionId, teamId, week string) (*models.Team, error) {
	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/team/%s/stats;type=week;week=%s", teamId, week)

	teamWeeklyResponse, err := AuthHttpXMLRequest(sessionId, url)
	if err != nil {
		return nil, err
	}

	teamWeekXMLResponse, err := MapToTeamWeek(teamWeeklyResponse)
	if err != nil {
		return nil, err
	}

	return teamWeekXMLResponse, nil
}

func GetPlayerStats(sessionId, playerId string) (*models.Player, error) {
	// Check if the player stats already exist in the database
	existingPlayer, err := repositories.GetPlayerStatsDB(playerId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to query existing player stats: %w", err)
	}

	// Check if stats exist and were recently updated
	if existingPlayer != nil {
		now := time.Now()
		if now.Before(existingPlayer.NextUpdate) {
			// Cached stats are still valid
			return existingPlayer, nil
		}
	}

	// If no recent stats or update needed, fetch from the Yahoo API
	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/player/%s/stats;type=season", playerId)
	playerStatsResponse, err := AuthHttpXMLRequest(sessionId, url)
	if err != nil {
		return nil, err
	}

	player, err := MapPlayer(playerStatsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to map player data: %w", err)
	}

	// Save the updated stats to the database
	err = repositories.SavePlayerStatsDB(*player)
	if err != nil {
		return nil, fmt.Errorf("failed to save player stats to the database: %w", err)
	}

	return player, nil
}

func GetPlayerRankLeague(sessionId, leagueId, playerId string) ([]models.PlayerRank, error) {
	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/leagues;league_keys=%s/players;player_keys=%s;out=ranks;ranks=season,last30days,last14days,last7days,projected_next7days,projected_next14days,projected_season_remaining", leagueId, playerId)
	playerRankResponse, err := AuthHttpXMLRequest(sessionId, url)
	if err != nil {
		return nil, err
	}

	playerRanks, err := MapToRank(playerRankResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to map player data: %w", err)
	}

	return playerRanks, nil
}

func GetAllNhlPlayersYahoo(sessionId string) ([]*models.YahooPlayer, error) {
	gameKey := "453"
	var allPlayers []*models.YahooPlayer
	start := 0
	count := 25

	for {
		url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/game/%s/players?start=%d&count=%d", gameKey, start, count)
		playersResponse, err := AuthHttpXMLRequest(sessionId, url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Yahoo players: %v", err)
		}

		players, err := MapToYahooPlayer(playersResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Yahoo XML to map: %v", err)
		}

		var playerPtrs []*models.YahooPlayer
		for i := range players {
			playerPtrs = append(playerPtrs, &players[i])
		}

		// If no players are found, break the loop
		if len(players) == 0 {
			break
		}

		// Append all players to result
		allPlayers = append(allPlayers, playerPtrs...)

		// Increment pagination
		start += count
		time.Sleep(500 * time.Millisecond)
	}

	err := repositories.SaveYahooPlayerToDB(allPlayers)
	if err != nil {
		return nil, fmt.Errorf("error saving players to DB: %w", err)
	}

	return allPlayers, nil
}

func GetAllTeamsInLeague(sessionId, leagueId string) ([]models.LeagueTeam, error) {

	leagueTeamsFromDB, err := repositories.GetAllLeagueTeamsFromDB(leagueId)
	if err != nil {
		log.Printf("Failed to get league teams from DB: %w", err)
	}

	if leagueTeamsFromDB != nil {
		if len(leagueTeamsFromDB) > 0 {
			return leagueTeamsFromDB, nil
		}
	}

	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/league/%s/teams", leagueId)
	leagueTeamResponse, err := AuthHttpXMLRequest(sessionId, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch players from league: %w", err)
	}

	teams, err := MapFantasyTeamsFromLeague(leagueTeamResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to convert api Response from league: %w", err)
	}

	err = repositories.SaveLeagueTeamsToDB(teams)
	if err != nil {
		log.Printf("Failed to save teams in DB: %w", err)
	}

	return teams, nil
}

type TeamMatchupResponse struct {
	Matchups    []*models.Matchup                 `json:"matchups"`
	TeamStats   []*models.TeamWeeklyStats         `json:"teamStats"`
	StatWinners []*models.StatWinnerWeeklyMatchup `json:"statWinners"`
}

func GetFTeamMatchups(sessionId, teamId string) (*TeamMatchupResponse, error) {
	url := fmt.Sprintf("https://fantasysports.yahooapis.com/fantasy/v2/team/%s/matchups", teamId)
	teamMatchupResponse, err := AuthHttpXMLRequest(sessionId, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch matchups for team: %w", err)
	}

	matchups, teamStats, statWinners, err := MapTeamMatchups(teamMatchupResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to map matchups for team: %w", err)
	}

	if err := repositories.SaveTeamMatchups(matchups, teamStats, statWinners); err != nil {
		return nil, fmt.Errorf("failed to save matchups to DB: %w", err)
	}

	return &TeamMatchupResponse{
		Matchups:    matchups,
		TeamStats:   teamStats,
		StatWinners: statWinners,
	}, nil
}
