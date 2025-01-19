package services

import (
	"fmt"
)

func ExtractLeaguesFromResponse(response map[string]interface{}) (map[string]interface{}, error) {
	leagues := make([]interface{}, 0) // Use []interface{} for JSON compatibility

	fantasyContent, ok := response["fantasy_content"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid fantasy_content")
	}

	users, ok := fantasyContent["users"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid users")
	}

	user, ok := users["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid user")
	}

	games, ok := user["games"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid games")
	}

	gameList, ok := games["game"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid game list")
	}

	// Iterate over each game and extract leagues
	for _, game := range gameList {
		gameMap, ok := game.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid game entry")
		}

		leaguesData, ok := gameMap["leagues"].(map[string]interface{})
		if !ok {
			continue // Skip games without leagues
		}

		leagueEntries, ok := leaguesData["league"].([]interface{})
		if ok {
			// Multiple leagues
			leagues = append(leagues, leagueEntries...)
		} else if leagueEntry, ok := leaguesData["league"].(map[string]interface{}); ok {
			// Single league
			leagues = append(leagues, leagueEntry)
		}
	}

	return map[string]interface{}{"leagues": leagues}, nil
}

func ExtractLeagueFromResponse(response map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
