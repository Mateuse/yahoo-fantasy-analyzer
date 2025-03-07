package services

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func MapToLeague(response map[string]interface{}) (*models.League, error) {
	// Extract the "league" map from the response
	leagueData, ok := response["league"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response: missing league data")
	}

	// Create a new League object
	league := &models.League{}

	// Map values to the League struct
	if val, ok := leagueData["league_id"].(string); ok {
		league.LeagueID = val
	}
	if val, ok := leagueData["league_key"].(string); ok {
		league.LeagueKey = val
	}
	if val, ok := leagueData["name"].(string); ok {
		league.Name = val
	}
	if val, ok := leagueData["url"].(string); ok {
		league.URL = val
	}
	if val, ok := leagueData["logo_url"].(string); ok {
		league.LogoURL = val
	}
	if val, ok := leagueData["draft_status"].(string); ok {
		league.DraftStatus = val
	}
	if val, ok := leagueData["num_teams"].(string); ok {
		league.NumTeams, _ = strconv.Atoi(val)
	}
	if val, ok := leagueData["scoring_type"].(string); ok {
		league.ScoringType = val
	}
	if val, ok := leagueData["league_type"].(string); ok {
		league.LeagueType = val
	}
	if val, ok := leagueData["felo_tier"].(string); ok {
		league.FeloTier = val
	}
	if val, ok := leagueData["current_week"].(string); ok {
		league.CurrentWeek, _ = strconv.Atoi(val)
	}
	if val, ok := leagueData["start_week"].(string); ok {
		league.StartWeek, _ = strconv.Atoi(val)
	}
	if val, ok := leagueData["start_date"].(string); ok {
		league.StartDate = utils.ParseDate(val)
	}
	if val, ok := leagueData["end_week"].(string); ok {
		league.EndWeek, _ = strconv.Atoi(val)
	}
	if val, ok := leagueData["end_date"].(string); ok {
		league.EndDate = utils.ParseDate(val)
	}
	if val, ok := leagueData["game_code"].(string); ok {
		league.GameCode = val
	}
	if val, ok := leagueData["season"].(string); ok {
		league.Season = val
	}
	if val, ok := leagueData["weekly_deadline"].(string); ok {
		league.WeeklyDeadline = val
	}
	if val, ok := leagueData["allow_add_to_dl_extra_pos"].(string); ok {
		league.AllowAddToDLExtraPos = val == "1"
	}
	if val, ok := leagueData["is_pro_league"].(string); ok {
		league.IsProLeague = val == "1"
	}
	if val, ok := leagueData["is_cash_league"].(string); ok {
		league.IsCashLeague = val == "1"
	}
	if val, ok := leagueData["is_plus_league"].(string); ok {
		league.IsPlusLeague = val == "1"
	}
	if val, ok := leagueData["league_update_timestamp"].(string); ok {
		league.LeagueUpdateTimestamp, _ = strconv.ParseInt(val, 10, 64)
	}

	statNames := make(map[string]string)
	// Map roster positions and stat modifiers
	if settings, ok := leagueData["settings"].(map[string]interface{}); ok {
		// Map roster positions
		if rosterPositions, ok := settings["roster_positions"].(map[string]interface{}); ok {
			if positions, ok := rosterPositions["roster_position"].([]interface{}); ok {
				for _, positionData := range positions {
					if pos, ok := positionData.(map[string]interface{}); ok {
						rosterPosition := models.RosterPosition{}
						if val, ok := pos["position"].(string); ok {
							rosterPosition.Position = val
						}
						if val, ok := pos["position_type"].(string); ok {
							rosterPosition.PositionType = val
						}
						if val, ok := pos["count"].(string); ok {
							rosterPosition.Count, _ = strconv.Atoi(val)
						}
						if val, ok := pos["is_starting_position"].(string); ok {
							rosterPosition.IsStartingPosition = val == "1"
						}
						league.RosterPositions = append(league.RosterPositions, rosterPosition)
					}
				}
			}
		}

		// Extract stat categories to populate statNames
		if statCategories, ok := settings["stat_categories"].(map[string]interface{}); ok {
			if stats, ok := statCategories["stats"].(map[string]interface{}); ok {
				if statList, ok := stats["stat"].([]interface{}); ok {
					for _, statData := range statList {
						if stat, ok := statData.(map[string]interface{}); ok {
							if statID, ok := stat["stat_id"].(string); ok {
								if statName, ok := stat["name"].(string); ok {
									statNames[statID] = statName
								}
							}
						}
					}
				}
			}
		}

		// Map stat modifiers and populate StatName
		if statModifiers, ok := settings["stat_modifiers"].(map[string]interface{}); ok {
			if stats, ok := statModifiers["stats"].(map[string]interface{}); ok {
				if statList, ok := stats["stat"].([]interface{}); ok {
					for _, statData := range statList {
						if stat, ok := statData.(map[string]interface{}); ok {
							statModifier := models.StatModifier{}
							if val, ok := stat["stat_id"].(string); ok {
								statModifier.StatID = val
								// Populate StatName using the statNames map
								if name, exists := statNames[val]; exists {
									statModifier.StatName = name
								}
							}
							if val, ok := stat["value"].(string); ok {
								statModifier.Value, _ = strconv.ParseFloat(val, 64)
							}
							league.StatModifiers = append(league.StatModifiers, statModifier)
						}
					}
				}
			}
		}
	}

	return league, nil
}

func MapToTeamWeek(response map[string]interface{}) (*models.Team, error) {
	teamData, ok := response["team"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response: missing team data")
	}

	team := &models.Team{}

	// Map values to the Team struct
	if val, ok := teamData["team_key"].(string); ok {
		team.TeamKey = val
	}
	if val, ok := teamData["team_id"].(string); ok {
		team.TeamID = val
	}
	if val, ok := teamData["name"].(string); ok {
		team.Name = val
	}
	if val, ok := teamData["url"].(string); ok {
		team.URL = val
	}

	// Handle nested team logos
	if logos, ok := teamData["team_logos"].([]interface{}); ok && len(logos) > 0 {
		if logoData, ok := logos[0].(map[string]interface{}); ok {
			if url, ok := logoData["url"].(string); ok {
				team.LogoURL = url
			}
		}
	}

	if val, ok := teamData["waiver_priority"].(string); ok {
		team.WaiverPriority, _ = strconv.Atoi(val)
	}
	if val, ok := teamData["number_of_moves"].(string); ok {
		team.NumberOfMoves, _ = strconv.Atoi(val)
	}
	if val, ok := teamData["number_of_trades"].(string); ok {
		team.NumberOfTrades, _ = strconv.Atoi(val)
	}
	if val, ok := teamData["league_scoring_type"].(string); ok {
		team.LeagueScoringType = val
	}
	if val, ok := teamData["draft_position"].(string); ok {
		team.DraftPosition, _ = strconv.Atoi(val)
	}

	// Map projected points
	if projectedPoints, ok := teamData["team_projected_points"].(map[string]interface{}); ok {
		if total, ok := projectedPoints["total"].(string); ok {
			team.ProjectedPoints = total
		}
	}

	// Map live projected points
	if liveProjectedPoints, ok := teamData["team_live_projected_points"].(map[string]interface{}); ok {
		if total, ok := liveProjectedPoints["total"].(string); ok {
			team.FinalPoints = total
		}
	}

	// Handle nested remaining games data
	if remainingGames, ok := teamData["team_remaining_games"].(map[string]interface{}); ok {
		if total, ok := remainingGames["total"].(map[string]interface{}); ok {
			if remaining, ok := total["remaining_games"].(string); ok {
				team.RemainingGames, _ = strconv.Atoi(remaining)
			}
			if completed, ok := total["completed_games"].(string); ok {
				team.CompletedGames, _ = strconv.Atoi(completed)
			}
		}
	}

	return team, nil
}

func MapPlayer(response map[string]interface{}) (*models.Player, error) {
	// Extract player data
	playerData, ok := response["player"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse player data")
	}

	// Extract name data
	nameData, ok := playerData["name"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse name data")
	}

	// Extract stats
	statsData := extractStats(playerData, "player_stats")
	advancedStatsData := extractStats(playerData, "player_advanced_stats")

	// Map the player object
	player := &models.Player{
		PlayerID:  utils.GetString(playerData, "player_id"),
		PlayerKey: utils.GetString(playerData, "player_key"),
		Name: models.PlayerName{
			FullName:   utils.GetString(nameData, "full"),
			FirstName:  utils.GetString(nameData, "first"),
			LastName:   utils.GetString(nameData, "last"),
			AsciiFirst: utils.GetString(nameData, "ascii_first"),
			AsciiLast:  utils.GetString(nameData, "ascii_last"),
		},
		TeamFullName:       utils.GetString(playerData, "editorial_team_full_name"),
		TeamAbbreviation:   utils.GetString(playerData, "editorial_team_abbr"),
		TeamURL:            utils.GetString(playerData, "editorial_team_url"),
		UniformNumber:      utils.GetString(playerData, "uniform_number"),
		DisplayPosition:    utils.GetString(playerData, "display_position"),
		HeadshotURL:        utils.GetString(playerData["headshot"].(map[string]interface{}), "url"),
		ImageURL:           utils.GetString(playerData, "image_url"),
		IsUndroppable:      utils.GetBool(playerData, "is_undroppable"),
		PositionType:       utils.GetString(playerData, "position_type"),
		EligiblePositions:  extractPositions(playerData["eligible_positions"]),
		PlayerNotes:        utils.GetBool(playerData, "has_player_notes"),
		RecentNotes:        utils.GetBool(playerData, "has_recent_player_notes"),
		PlayerNotesUpdated: utils.GetInt(playerData, "player_notes_last_timestamp"),
		Stats:              statsData,
		AdvancedStats:      advancedStatsData,
	}

	return player, nil
}

func extractStats(data map[string]interface{}, key string) []models.Stat {
	stats := []models.Stat{}

	// Access the stats container
	statsContainer, ok := data[key].(map[string]interface{})
	if !ok {
		fmt.Println("Stats container missing or invalid")
		return stats
	}

	// Access the "stats" map
	statsData, ok := statsContainer["stats"].(map[string]interface{})
	if !ok {
		fmt.Println("Stats key missing or invalid")
		return stats
	}

	// Access the "stat" field
	statField, exists := statsData["stat"]
	if !exists {
		fmt.Println("Stat field missing")
		return stats
	}

	// Handle the "stat" field
	switch v := statField.(type) {
	case []interface{}: // Multiple stats
		for _, stat := range v {
			statMap, ok := stat.(map[string]interface{})
			if ok {
				stats = append(stats, models.Stat{
					StatID: utils.GetString(statMap, "stat_id"),
					Value:  utils.GetString(statMap, "value"),
				})
			}
		}
	case map[string]interface{}: // Single stat
		stats = append(stats, models.Stat{
			StatID: utils.GetString(v, "stat_id"),
			Value:  utils.GetString(v, "value"),
		})
	default:
		fmt.Println("Stat field has unexpected type")
	}

	return stats
}

func extractPositions(positionsData interface{}) []string {
	positions := []string{}
	positionsArray, ok := positionsData.([]interface{})
	if !ok {
		return positions
	}

	for _, position := range positionsArray {
		if posStr, ok := position.(string); ok {
			positions = append(positions, posStr)
		}
	}

	return positions
}

func MapToRank(data map[string]interface{}) ([]models.PlayerRank, error) {
	var playerRanks []models.PlayerRank

	if leagues, ok := data["leagues"].(map[string]interface{}); ok {
		if league, ok := leagues["league"].(map[string]interface{}); ok {
			if players, ok := league["players"].(map[string]interface{}); ok {
				if player, ok := players["player"].(map[string]interface{}); ok {
					if playerRanksXML, ok := player["player_ranks"].(map[string]interface{}); ok {
						if playerRankArray, ok := playerRanksXML["player_rank"].([]interface{}); ok {
							for _, rankData := range playerRankArray {
								rankMap, ok := rankData.(map[string]interface{})
								if !ok {
									continue
								}

								// Extract rank details
								rankType := utils.GetString(rankMap, "rank_type")
								rankValue := utils.GetInt(rankMap, "rank_value")
								rankSeason := utils.GetString(rankMap, "rank_season")

								// Append rank to the slice
								playerRanks = append(playerRanks, models.PlayerRank{
									RankType:   rankType,
									RankValue:  rankValue,
									RankSeason: rankSeason,
								})
							}
						}
					}
				}
			}
		}
	}

	return playerRanks, nil
}

func mapNHLPlayer(data interface{}) (*models.NHLPlayer, error) {
	playerMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse player data")
	}

	player := &models.NHLPlayer{
		ID:             int(playerMap["id"].(float64)),
		Headshot:       utils.GetString(playerMap, "headshot"),
		SweaterNumber:  utils.GetInt(playerMap, "sweaterNumber"),
		PositionCode:   utils.GetString(playerMap, "positionCode"),
		ShootsCatches:  utils.GetString(playerMap, "shootsCatches"),
		HeightInInches: utils.GetInt(playerMap, "heightInInches"),
		WeightInPounds: utils.GetInt(playerMap, "weightInPounds"),
		HeightInCM:     utils.GetInt(playerMap, "heightInCentimeters"),
		WeightInKG:     utils.GetInt(playerMap, "weightInKilograms"),
		BirthDate:      utils.GetString(playerMap, "birthDate"),
		BirthCountry:   utils.GetString(playerMap, "birthCountry"),
	}

	// Extract first name
	if firstNameMap, ok := playerMap["firstName"].(map[string]interface{}); ok {
		player.FirstName = utils.GetString(firstNameMap, "default")
	}

	// Extract last name
	if lastNameMap, ok := playerMap["lastName"].(map[string]interface{}); ok {
		player.LastName = utils.GetString(lastNameMap, "default")
	}

	// Extract birth city
	if birthCityMap, ok := playerMap["birthCity"].(map[string]interface{}); ok {
		player.BirthCity = utils.GetString(birthCityMap, "default")
	}

	// Extract birth state (optional)
	if birthStateMap, ok := playerMap["birthStateProvince"].(map[string]interface{}); ok {
		player.BirthState = utils.GetString(birthStateMap, "default")
	}

	return player, nil
}

func MapToYahooPlayer(data map[string]interface{}) ([]models.YahooPlayer, error) {
	var players []models.YahooPlayer

	// Extract "game" object
	game, ok := data["game"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid game data")
	}

	// Check if "players" key exists before trying to access it
	playersMapRaw, exists := game["players"]
	if !exists {
		// No players in the response, return an empty slice (not an error)
		return []models.YahooPlayer{}, nil
	}

	// Extract "players" array
	playersMap, ok := playersMapRaw.(map[string]interface{})
	if !ok {
		log.Println("Players Reached")
		return nil, nil
	}

	playersData, ok := playersMap["player"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid player data")
	}
	// Loop through all player entries
	for _, playerEntry := range playersData {
		playerMap, ok := playerEntry.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract required fields
		playerID := utils.GetString(playerMap, "player_key")
		nameMap, ok := playerMap["name"].(map[string]interface{})
		if !ok {
			continue
		}
		fullName := utils.GetString(nameMap, "full")

		teamName := utils.GetString(playerMap, "editorial_team_full_name")

		headshotMap, ok := playerMap["headshot"].(map[string]interface{})
		if !ok {
			continue
		}
		headshotURL := utils.GetString(headshotMap, "url")

		// Append to slice
		players = append(players, models.YahooPlayer{
			ID:          playerID,
			FullName:    fullName,
			TeamName:    teamName,
			HeadshotURL: headshotURL,
		})
	}

	return players, nil
}

func MapFantasyTeamsFromLeague(data map[string]interface{}) ([]models.LeagueTeam, error) {
	var leagueTeams []models.LeagueTeam

	// Access the "league" object
	league, ok := data["league"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing 'league' data")
	}

	// Extract league ID
	leagueID, ok := league["league_key"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'league_key' in league data")
	}

	// Access the "teams" object
	teamsData, ok := league["teams"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing 'teams' in league data")
	}

	// Extract "team" list
	teamList, ok := teamsData["team"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing 'team' list in teams data")
	}

	// Iterate over each team
	for _, teamEntry := range teamList {
		teamMap, ok := teamEntry.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract team fields
		teamID, _ := teamMap["team_key"].(string)
		teamName, _ := teamMap["name"].(string)

		// Extract team logo (if exists)
		logo := ""
		if logos, exists := teamMap["team_logos"].(map[string]interface{}); exists {
			if logoEntries, exists := logos["team_logo"].(map[string]interface{}); exists {
				logo, _ = logoEntries["url"].(string)
			}
		}

		// Create the LeagueTeam struct
		team := models.LeagueTeam{
			TeamId:   teamID,
			LeagueId: leagueID,
			Name:     teamName,
			Logo:     logo,
		}
		leagueTeams = append(leagueTeams, team)
	}

	return leagueTeams, nil
}

func MapNHLGameStat(data map[string]interface{}) (*models.PlayerGameStat, error) {
	// Ensure data is valid
	if data == nil {
		return nil, errors.New("input data is nil")
	}

	playerGameStat := &models.PlayerGameStat{}

	// Safely extract and convert each field
	if gameID, exists := data["gameId"].(float64); exists {
		playerGameStat.GameID = fmt.Sprintf("%.0f", gameID)
	}

	if teamAbbrev, exists := data["teamAbbrev"].(string); exists {
		playerGameStat.TeamAbbrev = teamAbbrev
	}

	if homeRoadFlag, exists := data["homeRoadFlag"].(string); exists {
		playerGameStat.HomeRoadFlag = homeRoadFlag
	}

	if gameDate, exists := data["gameDate"].(string); exists {
		playerGameStat.GameDate = gameDate
	}

	if goals, exists := data["goals"].(float64); exists {
		playerGameStat.Goals = int(goals)
	}

	if assists, exists := data["assists"].(float64); exists {
		playerGameStat.Assists = int(assists)
	}

	if points, exists := data["points"].(float64); exists {
		playerGameStat.Points = int(points)
	}

	if plusMinus, exists := data["plusMinus"].(float64); exists {
		playerGameStat.PlusMinus = int(plusMinus)
	}

	if powerPlayGoals, exists := data["powerPlayGoals"].(float64); exists {
		playerGameStat.PowerPlayGoals = int(powerPlayGoals)
	}

	if powerPlayPoints, exists := data["powerPlayPoints"].(float64); exists {
		playerGameStat.PowerPlayPoints = int(powerPlayPoints)
	}

	if gameWinningGoals, exists := data["gameWinningGoals"].(float64); exists {
		playerGameStat.GameWinningGoals = int(gameWinningGoals)
	}

	if otGoals, exists := data["otGoals"].(float64); exists {
		playerGameStat.OTGoals = int(otGoals)
	}

	if shots, exists := data["shots"].(float64); exists {
		playerGameStat.Shots = int(shots)
	}

	if shifts, exists := data["shifts"].(float64); exists {
		playerGameStat.Shifts = int(shifts)
	}

	if shorthandedGoals, exists := data["shorthandedGoals"].(float64); exists {
		playerGameStat.ShorthandedGoals = int(shorthandedGoals)
	}

	if shorthandedPoints, exists := data["shorthandedPoints"].(float64); exists {
		playerGameStat.ShorthandedPoints = int(shorthandedPoints)
	}

	if opponentAbbrev, exists := data["opponentAbbrev"].(string); exists {
		playerGameStat.OpponentAbbrev = opponentAbbrev
	}

	if pim, exists := data["pim"].(float64); exists {
		playerGameStat.PIM = int(pim)
	}

	if toi, exists := data["toi"].(string); exists {
		playerGameStat.TOI = toi
	}

	// Extract `team` from `commonName.default`
	if commonName, exists := data["commonName"].(map[string]interface{}); exists {
		if team, ok := commonName["default"].(string); ok {
			playerGameStat.Team = team
		}
	}

	// Extract `opponent` from `opponentCommonName.default`
	if opponentCommonName, exists := data["opponentCommonName"].(map[string]interface{}); exists {
		if opponent, ok := opponentCommonName["default"].(string); ok {
			playerGameStat.Opponent = opponent
		}
	}

	return playerGameStat, nil
}

func MapTeamMatchups(data map[string]interface{}) ([]*models.Matchup, []*models.TeamWeeklyStats, []*models.StatWinnerWeeklyMatchup, error) {
	teamData, ok := data["team"].(map[string]interface{})
	if !ok {
		return nil, nil, nil, errors.New("invalid team data")
	}
	matchupsData, ok := teamData["matchups"].(map[string]interface{})
	if !ok {
		return nil, nil, nil, errors.New("invalid matchup list format")
	}

	matchupList, ok := matchupsData["matchup"].([]interface{})
	if !ok {
		return nil, nil, nil, errors.New("invalid matchup list format")
	}

	var matchups []*models.Matchup
	var teamStats []*models.TeamWeeklyStats
	var statWinners []*models.StatWinnerWeeklyMatchup

	for _, item := range matchupList {
		matchupMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		week := utils.GetString(matchupMap, "week")
		winnerTeam := utils.GetString(matchupMap, "winner_team_key")
		if winnerTeam == "" {
			continue
		}

		teamsData, ok := matchupMap["teams"].(map[string]interface{})
		if !ok {
			continue
		}

		teams, ok := teamsData["team"].([]interface{})
		if !ok {
			continue
		}

		var teamKeys []string
		for _, team := range teams {
			teamMap, ok := team.(map[string]interface{})
			if !ok {
				continue
			}
			teamKey := utils.GetString(teamMap, "team_key")
			teamKeys = append(teamKeys, teamKey)

			teamStatsData := extractTeamStats(teamMap, teamKey, week)
			teamStats = append(teamStats, teamStatsData)
		}

		losingTeam := teamKeys[0]
		if teamKeys[0] == winnerTeam {
			losingTeam = teamKeys[1]
		}

		matchupId := utils.GenerateMatchupKey(winnerTeam, losingTeam, week)
		matchups = append(matchups, &models.Matchup{
			MatchupKey:  matchupId,
			Week:        week,
			WinningTeam: winnerTeam,
			LosingTeam:  losingTeam,
		})

		if statWinnersData, ok := matchupMap["stat_winners"].(map[string]interface{}); ok {
			extractedStatWinners := extractStatWinners(statWinnersData, week, matchupId)
			statWinners = append(statWinners, extractedStatWinners...)
		}
	}

	return matchups, teamStats, statWinners, nil
}

func extractTeamStats(teamMap map[string]interface{}, teamKey, week string) *models.TeamWeeklyStats {
	var points float64

	if pointsData, ok := teamMap["team_points"].(map[string]interface{}); ok {
		points = utils.GetFloat(pointsData, "total")
	}

	var statsSlice []models.Stat

	if statsData, ok := teamMap["team_stats"].(map[string]interface{}); ok {
		if statsList, ok := statsData["stats"].(map[string]interface{}); ok {
			if statArray, ok := statsList["stat"].([]interface{}); ok {
				for _, stat := range statArray {
					statMap, ok := stat.(map[string]interface{})
					if !ok {
						continue
					}
					statId := utils.GetString(statMap, "stat_id")
					value := utils.GetString(statMap, "value")

					statsSlice = append(statsSlice, models.Stat{
						StatID: statId,
						Value:  value,
					})
				}
			}
		}
	}

	return &models.TeamWeeklyStats{
		ID:      fmt.Sprintf("%s-%s", teamKey, week),
		TeamKey: teamKey,
		Week:    week,
		Stats:   statsSlice,
		Points:  points,
	}
}

func extractStatWinners(statWinnersData map[string]interface{}, week, matchupId string) []*models.StatWinnerWeeklyMatchup {
	var statWinners []*models.StatWinnerWeeklyMatchup

	if statWinnerList, ok := statWinnersData["stat_winner"].([]interface{}); ok {
		for _, stat := range statWinnerList {
			statMap, ok := stat.(map[string]interface{})
			if !ok {
				continue
			}

			statId := utils.GetString(statMap, "stat_id")
			isTied := false
			if tied, exists := statMap["is_tied"]; exists {
				isTied, _ = strconv.ParseBool(fmt.Sprintf("%v", tied))
			}

			winnerTeam := ""
			if winner, exists := statMap["winner_team_key"]; exists {
				winnerTeam = fmt.Sprintf("%v", winner)
			}

			statWinners = append(statWinners, &models.StatWinnerWeeklyMatchup{
				Week:           week,
				MatchupKey:     matchupId,
				StatID:         statId,
				WinningTeamKey: winnerTeam,
				IsTied:         isTied,
			})
		}
	}
	return statWinners
}
