package services

import (
	"fmt"
	"log"
	"time"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/repositories"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func SaveAllTeamsSchedule() error {
	teamAbbrs := utils.GetNHLTeamAbbreviations()

	for abbr := range teamAbbrs {
		err := GetTeamSchedule(abbr)
		if err != nil {
			log.Printf("failed to save schedule for team %s: %v", abbr, err)
			return err
		}
	}

	return nil
}

func GetTeamSchedule(abbr string) error {
	url := fmt.Sprintf("https://api-web.nhle.com/v1/club-schedule-season/%s/now", abbr)

	response, err := GetHttpRequest(url)
	if err != nil {
		return fmt.Errorf("failed to fetch schedule for team %s: %w", abbr, err)
	}

	// Parse the response into the required structure
	gamesData, ok := response["games"].([]interface{})
	if !ok {
		return fmt.Errorf("unexpected response format for team %s", abbr)
	}

	for _, gameData := range gamesData {
		gameMap, ok := gameData.(map[string]interface{})
		if !ok {
			continue // Skip invalid entries
		}

		startTimeStr := gameMap["startTimeUTC"].(string)
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return fmt.Errorf("failed to parse start time for game ID %v: %w", gameMap["id"], err)
		}

		schedule := models.ScheduleGame{
			ID:             int64(gameMap["id"].(float64)),
			Season:         int(gameMap["season"].(float64)),
			GameType:       int(gameMap["gameType"].(float64)),
			GameDate:       gameMap["gameDate"].(string),
			StartTimeUTC:   startTime, // Use parsed time
			HomeTeamAbbrev: gameMap["homeTeam"].(map[string]interface{})["abbrev"].(string),
			AwayTeamAbbrev: gameMap["awayTeam"].(map[string]interface{})["abbrev"].(string),
		}

		if err := repositories.SaveScheduleGameInDB(&schedule); err != nil {
			return fmt.Errorf("failed to save schedule for game ID %d: %w", schedule.ID, err)
		}
	}

	return nil
}

func GetTeamRoster(teamAbrev, season string) ([]*models.NHLPlayer, error) {
	url := fmt.Sprintf("https://api-web.nhle.com/v1/roster/%s/%s", teamAbrev, season)

	response, err := GetHttpRequest(url)
	if err != nil {
		return nil, err
	}

	// Extract players from response
	var players []*models.NHLPlayer

	// Extract forwards
	if forwards, ok := response["forwards"].([]interface{}); ok {
		for _, f := range forwards {
			player, err := mapNHLPlayer(f)
			player.Team = utils.GetNHLTeamAbbreviations()[teamAbrev]
			if err == nil {
				players = append(players, player)
			}
		}
	}

	// Extract defensemen
	if defensemen, ok := response["defensemen"].([]interface{}); ok {
		for _, d := range defensemen {
			player, err := mapNHLPlayer(d)
			player.Team = utils.GetNHLTeamAbbreviations()[teamAbrev]
			if err == nil {
				players = append(players, player)
			}
		}
	}

	// Extract goalies
	if goalies, ok := response["goalies"].([]interface{}); ok {
		for _, g := range goalies {
			player, err := mapNHLPlayer(g)
			player.Team = utils.GetNHLTeamAbbreviations()[teamAbrev]
			if err == nil {
				players = append(players, player)
			}
		}
	}

	err = repositories.SaveNhlPlayerToDB(players)
	if err != nil {
		return nil, fmt.Errorf("failed to save players: %w", err)
	}

	return players, nil
}

func MapNhlPlayerToYahoo() error {
	var mappings []models.PlayerIDMapping

	nhlPlayers, err := repositories.GetNhlPlayers()
	if err != nil {
		return err
	}
	nhlPlayerMap := make(map[string]int)
	for _, nhlPlayer := range nhlPlayers {
		fullName := nhlPlayer.FirstName + " " + nhlPlayer.LastName
		nhlPlayerMap[fullName] = nhlPlayer.ID
	}

	yahooPlayers, err := repositories.GetYahooPlayers()
	if err != nil {
		return err
	}

	for _, yahooPlayer := range yahooPlayers {
		if nhlID, found := nhlPlayerMap[yahooPlayer.FullName]; found {
			mappings = append(mappings, models.PlayerIDMapping{
				YahooPlayerID: yahooPlayer.ID,
				NHLPlayerID:   fmt.Sprintf("%d", nhlID),
				PlayerName:    yahooPlayer.FullName,
				TeamAbbr:      yahooPlayer.TeamName,
			})
		}
	}

	err = repositories.SavePlayerIDMappingToDB(mappings)
	if err != nil {
		return err
	}

	return nil
}

func GetPlayerGameStatsNHL(playerId, season string) ([]*models.PlayerGameStat, error) {
	url := fmt.Sprintf("https://api-web.nhle.com/v1/player/%s/game-log/%s/2", playerId, season)

	response, err := GetHttpRequest(url)
	if err != nil {
		return nil, err
	}

	var playerGameStats []*models.PlayerGameStat

	if gameLog, exists := response["gameLog"].([]interface{}); exists {
		for _, game := range gameLog {
			game, err := MapNHLGameStat(game.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			game.PlayerID = playerId
			playerGameStats = append(playerGameStats, game)
		}
	}
	err = repositories.SavePlayerGameStats(playerGameStats)
	if err != nil {
		return nil, err
	}

	return playerGameStats, nil
}
