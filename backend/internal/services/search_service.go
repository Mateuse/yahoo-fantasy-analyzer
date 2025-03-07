package services

import (
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/repositories"
)

func GetPlayerByName(userSession, playerName string) (*models.PlayerDetails, error) {

	playerIds, err := repositories.GetMappedPlayerByName(playerName)
	if err != nil {
		return nil, err
	}

	yahooPlayerId := playerIds.YahooPlayerID
	nhlPlayerId := playerIds.NHLPlayerID

	yahooPlayer, err := GetPlayerStats(userSession, yahooPlayerId)
	if err != nil {
		return nil, err
	}

	nhlPlayer, err := repositories.GetNhlPlayerById(nhlPlayerId)
	if err != nil {
		return nil, err
	}

	player := &models.PlayerDetails{
		Player:            *nhlPlayer,
		YahooPlayerID:     yahooPlayerId,
		EligiblePositions: yahooPlayer.EligiblePositions,
		Stats:             yahooPlayer.Stats,
		AdvancedStats:     yahooPlayer.AdvancedStats,
	}

	return player, nil
}
